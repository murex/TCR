/*
Copyright (c) 2023 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package p4

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/utils"
	"github.com/murex/tcr/vcs"
	"github.com/murex/tcr/vcs/shell"
	"github.com/spf13/afero"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"path/filepath"
	"strconv"
	"strings"
)

// Name provides the name for this VCS implementation
const Name = "p4"

// p4Impl provides the implementation of the Perforce interface
type p4Impl struct {
	baseDir              string
	rootDir              string
	clientName           string
	filesystem           afero.Fs
	runP4Function        func(params ...string) (output []byte, err error)
	traceP4Function      func(params ...string) (err error)
	runPipedP4Function   func(toCmd shell.Command, params ...string) (output []byte, err error)
	tracePipedP4Function func(toCmd shell.Command, params ...string) (err error)
}

// New initializes the p4 implementation based on the provided directory from local clone
func New(dir string) (vcs.Interface, error) {
	return newP4Impl(plainOpen, dir, false)
}

func newP4Impl(initDepotFs func() afero.Fs, dir string, testFlag bool) (*p4Impl, error) {
	var p = p4Impl{
		baseDir:              dir,
		filesystem:           initDepotFs(),
		runP4Function:        runP4Command,
		traceP4Function:      traceP4Command,
		runPipedP4Function:   runPipedP4Command,
		tracePipedP4Function: tracePipedP4Command,
	}

	if testFlag {
		// For test purpose only: tests should run and pass without having p4 installed and with no p4 server available
		p.clientName = "test"
		p.rootDir = dir
	} else {
		p.clientName = GetP4ClientName()
		var err error
		p.rootDir, err = GetP4RootDir()
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsSubPathOf(p.baseDir, p.GetRootDir()) {
		return nil, fmt.Errorf("directory %s does not belong to a p4 depot", p.baseDir)
	}
	return &p, nil
}

// Name returns VCS name
func (*p4Impl) Name() string {
	return Name
}

// SessionSummary provides a short description related to current VCS session summary
func (p *p4Impl) SessionSummary() string {
	return fmt.Sprintf("%s client \"%s\"", p.Name(), p.clientName)
}

// plainOpen is the regular function used to open a p4 depot
func plainOpen() afero.Fs {
	return afero.NewOsFs()
}

// GetRootDir returns the root directory path
func (p *p4Impl) GetRootDir() string {
	return p.rootDir
}

// GetRemoteName returns the current p4 "remote name"
func (*p4Impl) GetRemoteName() string {
	// Always return an empty string
	return ""
}

// IsRemoteEnabled indicates if p4 remote operations are enabled
// remote is always enabled with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) IsRemoteEnabled() bool {
	return true
}

// GetWorkingBranch returns the current p4 working branch
func (*p4Impl) GetWorkingBranch() string {
	// For now, always return an empty string
	return ""
}

// IsOnRootBranch indicates if p4 is currently on its root branch or not.
func (*p4Impl) IsOnRootBranch() bool {
	// For now, always return false
	return false
}

// Add adds the listed paths to p4 changelist.
func (p *p4Impl) Add(paths ...string) error {
	return p.reconcile(paths...)
}

func (p *p4Impl) reconcile(paths ...string) error {
	p4Args := []string{"reconcile", "-a", "-e", "-d"}
	if len(paths) == 0 {
		p4Args = append(p4Args, filepath.Join(p.baseDir, "/..."))
	} else {
		p4Args = append(p4Args, paths...)
	}
	return p.traceP4(p4Args...)
}

type changeList struct {
	number string
}

// Commit commits changes to p4 index.
// With current implementation, "amend" parameter is ignored.
func (p *p4Impl) Commit(messages ...string) error {
	cl, err := p.createChangeList(messages...)
	if err != nil {
		report.PostError(err)
		return err
	}
	return p.submitChangeList(cl)
}

// RevertLocal restores to last commit for the provided path.
func (p *p4Impl) RevertLocal(path string) error {
	// in order to work, p4 revert requires that the file be reconciled beforehand
	err := p.reconcile(path)
	if err != nil {
		report.PostWarning(err)
	}
	// Command: p4 revert <path>
	return p.traceP4("revert", path)
}

// Revert runs a p4 revert operation.
// TODO: VCS Revert - p4 revert
func (*p4Impl) RollbackLastCommit() error {
	return errors.New("VCS revert operation not yet available for p4")
}

// Push runs a push operation.
func (*p4Impl) Push() error {
	// Nothing to do in case of p4, as the "p4 submit" done in Commit()
	// already "pushed" the changes to the server.
	return nil
}

// Pull runs a pull operation ("p4 sync")
func (p *p4Impl) Pull() error {
	path, err := p.toP4ClientPath(p.baseDir)
	if err != nil {
		return err
	}
	return p.traceP4("sync", path)
}

// Diff returns the list of files modified since last commit with diff info for each file
func (p *p4Impl) Diff() (diffs vcs.FileDiffs, err error) {
	var p4Output []byte
	p4Output, err = p.runP4("diff", "-f", "-Od", "-dl", "-ds", filepath.Join(p.baseDir, "/..."))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(p4Output))
	var filename string
	var added, removed int
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		switch fields[0] {
		case "====":
			// ==== <depot_path>>#<revision_number> - <local_path> ====
			filename = filepath.Clean(fields[3])
		case "add":
			// add <x> chunks <y> lines
			added, _ = strconv.Atoi(fields[3])
		case "deleted":
			// deleted <x> chunks <y> lines
			removed, _ = strconv.Atoi(fields[3])
		case "changed":
			// changed <x> chunks <y_before> / <y_after> lines
			changedBefore, _ := strconv.Atoi(fields[3])
			changedAfter, _ := strconv.Atoi(fields[5])
			diffs = append(diffs, vcs.NewFileDiff(filename, added+changedAfter, removed+changedBefore))
		default:
			return nil, fmt.Errorf("unrecognized p4 diff output: %s", scanner.Text())
		}
	}
	return diffs, nil
}

// Log returns the list of p4 log items compliant with the provided msgFilter.
// When no msgFilter is provided, returns all p4 log items unfiltered.
// TODO: VCS Log - p4 changes ./...
func (*p4Impl) Log(_ func(msg string) bool) (logs vcs.LogItems, err error) {
	return nil, errors.New("VCS log operation not yet available for p4")
}

// EnableAutoPush sets a flag allowing to turn on/off p4 auto-push operations.
// Auto-push is always on with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) EnableAutoPush(_ bool) {
	// nothing to do here.
}

// IsAutoPushEnabled indicates if p4 auto-push operations are turned on.
// Auto-push is always on with p4 due its architecture (all changes occur directly on the server)
func (*p4Impl) IsAutoPushEnabled() bool {
	return true
}

// CheckRemoteAccess returns true if p4 remote can be accessed.
func (*p4Impl) CheckRemoteAccess() bool {
	return true
}

// traceP4 runs a p4 command and traces its output.
// The command is launched from the p4 root directory
func (p *p4Impl) traceP4(args ...string) error {
	return p.traceP4Function(p.buildP4Args(args...)...)
}

// runP4 calls p4 command in a separate process and returns its output traces
// The command is launched from the p4 root directory
func (p *p4Impl) runP4(args ...string) (output []byte, err error) {
	return p.runP4Function(p.buildP4Args(args...)...)
}

func (p *p4Impl) runPipedP4(toCmd shell.Command, args ...string) (output []byte, err error) {
	return p.runPipedP4Function(toCmd, p.buildP4Args(args...)...)
}

func (p *p4Impl) buildP4Args(args ...string) []string {
	return append([]string{"-d", p.GetRootDir(), "-c", p.clientName}, args...)
}

func (p *p4Impl) createChangeList(messages ...string) (*changeList, error) {
	// Command: p4 --field "Description=<message>" change -o | p4 change -i
	out, err := p.runPipedP4(newP4Command(p.buildP4Args("change", "-i")...),
		"-Q", "utf8",
		"--field", buildDescriptionField(shell.GetAttributes(), messages...),
		"change", "-o")
	if err != nil {
		return nil, err
	}
	// Output: "Change <change list number> created ..."
	words := strings.Split(string(out), " ")
	if len(words) < 3 {
		return nil, fmt.Errorf("unexpected p4 change trace: %s", out)
	}
	clNumber := strings.Split(string(out), " ")[1]
	return &changeList{clNumber}, err
}

func buildDescriptionField(attr shell.Attributes, messages ...string) string {
	var builder strings.Builder
	_, _ = builder.WriteString("Description=")
	for _, message := range messages {
		_, _ = builder.WriteString(convertLine(attr.Encoding, message))
		_, _ = builder.WriteString(attr.EOL)
	}
	return builder.String()
}

func convertLine(charMap *charmap.Charmap, message string) string {
	if charMap == nil {
		// By default, we use UTF-8, which is the default with Go strings
		return message
	}
	var b bytes.Buffer
	converter := transform.NewWriter(&b, charMap.NewDecoder())
	_, _ = converter.Write([]byte(message))
	_ = converter.Close()
	return b.String()
}

func (p *p4Impl) submitChangeList(cl *changeList) error {
	if cl == nil {
		report.PostWarning("Empty changelist!")
		return errors.New("empty p4 changelist")
	}
	return p.traceP4("submit", "-c", cl.number)
}

func (p *p4Impl) toP4ClientPath(dir string) (string, error) {
	cleanDir := filepath.Clean(dir)
	cleanRoot := filepath.Clean(p.rootDir)

	if dir == "" {
		return "", errors.New("can not convert an empty path")
	}
	if !utils.IsSubPathOf(cleanDir, cleanRoot) {
		return "", errors.New("path is outside p4 root directory")
	}

	relativePath := strings.Replace(cleanDir, cleanRoot, "", 1)
	slashedPath := strings.ReplaceAll(relativePath, "\\", "/")
	if !strings.HasPrefix(slashedPath, "/") {
		slashedPath = "/" + slashedPath
	}
	if !strings.HasSuffix(slashedPath, "/") {
		slashedPath = slashedPath + "/"
	}
	return "//" + p.clientName + slashedPath + "...", nil
}
