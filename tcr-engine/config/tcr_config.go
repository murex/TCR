/*
Copyright (c) 2022 Murex

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

package config

import (
	"fmt"
	"github.com/murex/tcr/tcr-engine/params"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// TcrConfig wraps all possible TCR configuration parameters
type TcrConfig struct {
	BaseDir          *StringParam
	WorkDir          *StringParam
	ConfigDir        *StringParam
	Language         *StringParam
	Toolchain        *StringParam
	PollingPeriod    *DurationParam
	MobTimerDuration *DurationParam
	AutoPush         *BoolParam
	CommitFailures   *BoolParam
}

func (c TcrConfig) reset() {
	c.Language.reset()
	c.Toolchain.reset()
	c.PollingPeriod.reset()
	c.MobTimerDuration.reset()
	c.AutoPush.reset()
	c.CommitFailures.reset()
}

// Config is the placeholder for all TCR configuration parameters
var Config = TcrConfig{}

const (
	configDirRoot  = ".tcr"
	configFileBase = "config"
	configFileType = "yml"
	configFileName = configFileBase + "." + configFileType
)

var (
	// DirPathGetter is the function allowing to retrieve the configuration directory path. It should
	// not be changed unless when running tests.
	DirPathGetter     = GetConfigDirPath
	configDirPath     string
	configTraceWriter io.Writer
)

// StandardInit initializes TCR configuration
func StandardInit() {
	// With standard initialization, configuration trace is sent to stderr
	initConfig(os.Stderr)
}

func initConfig(writer io.Writer) {
	setConfigTrace(writer)
	// Initialize configuration directory path
	initConfigDirPath()
	// Make sure configuration directory exists
	createConfigDir()

	initTcrConfig()
	initToolchainConfig()
	initLanguageConfig()
}

func initTcrConfig() {
	// Viper setup
	configFilePath := filepath.Join(configDirPath, configFileName)
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configFileType)
	viper.SetConfigName(configFileName)
	viper.SetConfigFile(configFilePath)
	viper.SetEnvPrefix(settings.ApplicationName)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		trace("Loading configuration: ", viper.ConfigFileUsed())
	} else {
		trace("No configuration file found")
	}
}

func setConfigTrace(w io.Writer) {
	if w != nil {
		configTraceWriter = w
	}
}

func initConfigDirPath() {
	if Config.ConfigDir == nil || Config.ConfigDir.GetValue() == "" {
		// If configuration directory is not specified, we use by default the current directory
		configDirPath = filepath.Join(".", configDirRoot)
	} else {
		configDirPath = filepath.Join(Config.ConfigDir.GetValue(), configDirRoot)
	}
}

// GetConfigDirPath returns the configuration directory path
func GetConfigDirPath() string {
	return configDirPath
}

func createConfigDir() {
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		trace("Creating TCR configuration directory: ", configDirPath)
		err := os.MkdirAll(configDirPath, os.ModePerm)
		if err != nil {
			trace("Error creating TCR configuration directory: ", err)
		}
	}
}

// Save saves TCR configuration
func Save() {
	createConfigDir()
	saveTcrConfig()
	saveToolchainConfigs()
	saveLanguageConfigs()
}

func saveTcrConfig() {
	trace("Saving configuration: ", viper.ConfigFileUsed())
	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			_ = viper.WriteConfigAs(viper.ConfigFileUsed())
		} else {
			trace("Error while saving configuration file: ", err)
		}
	}
}

// Reset resets TCR configuration to default value
func Reset() {
	trace("Resetting configuration to default values")
	resetTcrConfig()
	resetToolchainConfigs()
	resetLanguageConfigs()
	Save()
}

func resetTcrConfig() {
	Config.reset()
}

// Show displays current TCR configuration
func Show() {
	trace()
	showTrcConfig()
	showToolchainConfigs()
	showLanguageConfigs()
}

func showTrcConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	trace("TCR configuration:")
	for _, key := range keys {
		showConfigValue(key, viper.Get(key))
	}
}

func showConfigValue(key string, value interface{}) {
	trace("- ", key, ": ", value)
}

func trace(a ...interface{}) {
	if configTraceWriter != nil {
		_, _ = fmt.Fprintln(configTraceWriter, "["+settings.ApplicationName+"]", fmt.Sprint(a...))
	}
}

// AddParameters adds parameter to the provided command cmd
func AddParameters(cmd *cobra.Command, defaultDir string) {
	Config.BaseDir = AddBaseDirParamWithDefault(cmd, defaultDir)
	Config.WorkDir = AddWorkDirParamWithDefault(cmd, defaultDir)
	Config.ConfigDir = AddConfigDirParam(cmd)
	Config.Language = AddLanguageParam(cmd)
	Config.Toolchain = AddToolchainParam(cmd)
	Config.PollingPeriod = AddPollingPeriodParam(cmd)
	Config.MobTimerDuration = AddMobTimerDurationParam(cmd)
	Config.AutoPush = AddAutoPushParam(cmd)
	Config.CommitFailures = AddCommitFailuresParam(cmd)
}

// UpdateEngineParams updates TCR engine parameters based on configuration values
func UpdateEngineParams(params *params.Params) {
	params.BaseDir = Config.BaseDir.GetValue()
	params.WorkDir = Config.WorkDir.GetValue()
	params.ConfigDir = Config.ConfigDir.GetValue()
	params.MobTurnDuration = Config.MobTimerDuration.GetValue()
	params.Language = Config.Language.GetValue()
	params.Toolchain = Config.Toolchain.GetValue()
	params.PollingPeriod = Config.PollingPeriod.GetValue()
	params.AutoPush = Config.AutoPush.GetValue()
	params.CommitFailures = Config.CommitFailures.GetValue()
}
