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

package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type sessionInfo struct {
	BaseDir           string `json:"baseDir"`
	WorkDir           string `json:"workDir"`
	LanguageName      string `json:"language"`
	ToolchainName     string `json:"toolchain"`
	VCSName           string `json:"vcsName"`
	VCSSessionSummary string `json:"vcsSession"`
	CommitOnFail      bool   `json:"commitOnFail"`
	Flavor            string `json:"flavor"`
	GitAutoPush       bool   `json:"gitAutoPush"`
	MessageSuffix     string `json:"messageSuffix"`
}

// SessionInfoGetHandler handles HTTP GET requests on TCR settings information
func SessionInfoGetHandler(c *gin.Context) {
	tcr := getTCRInstance(c)
	info := tcr.GetSessionInfo()
	data := sessionInfo{
		BaseDir:           info.BaseDir,
		WorkDir:           info.WorkDir,
		LanguageName:      info.LanguageName,
		ToolchainName:     info.ToolchainName,
		VCSName:           info.VCSName,
		VCSSessionSummary: info.VCSSessionSummary,
		CommitOnFail:      info.CommitOnFail,
		Flavor:            info.Flavor,
		GitAutoPush:       info.GitAutoPush,
		MessageSuffix:     info.MessageSuffix,
	}
	c.IndentedJSON(http.StatusOK, data)
}
