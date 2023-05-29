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

package config

import (
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/toolchain"
	"github.com/murex/tcr/utils"
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
	VCS              *StringParam
	MessageSuffix    *StringParam
	Trace            *StringParam
}

func (c TcrConfig) reset() {
	c.Language.reset()
	c.Toolchain.reset()
	c.PollingPeriod.reset()
	c.MobTimerDuration.reset()
	c.AutoPush.reset()
	c.CommitFailures.reset()
	c.VCS.reset()
	c.MessageSuffix.reset()
	c.Trace.reset()
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
	configDirPath string
)

// StandardInit initializes TCR configuration
func StandardInit() {
	// With standard initialization, configuration Trace is sent to stderr
	initConfig(os.Stderr)
}

func initConfig(writer io.Writer) {
	utils.SetSimpleTrace(writer)
	// Initialize configuration directory path
	initConfigDirPath()
	// Make sure configuration directory exists
	createConfigDir()

	initTCRConfig()
	toolchain.InitConfig(configDirPath)
	language.InitConfig(configDirPath)
}

func initTCRConfig() {
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
		utils.Trace("Loading configuration: ", viper.ConfigFileUsed())
	} else {
		utils.Trace("No configuration file found")
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
		utils.Trace("Creating TCR configuration directory: ", configDirPath)
		err := os.MkdirAll(configDirPath, os.ModePerm)
		if err != nil {
			utils.Trace("Error creating TCR configuration directory: ", err)
		}
	}
}

// Save saves TCR configuration
func Save() {
	createConfigDir()
	saveTCRConfig()
	toolchain.SaveConfigs()
	language.SaveConfigs()
}

func saveTCRConfig() {
	utils.Trace("Saving configuration: ", viper.ConfigFileUsed())
	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			_ = viper.WriteConfigAs(viper.ConfigFileUsed())
		} else {
			utils.Trace("Error while saving configuration file: ", err)
		}
	}
}

// Reset resets TCR configuration to default value
func Reset() {
	utils.Trace("Resetting configuration to default values")
	resetTCRConfig()
	toolchain.ResetConfigs()
	language.ResetConfigs()
	Save()
}

func resetTCRConfig() {
	Config.reset()
}

// Show displays current TCR configuration
func Show() {
	utils.Trace()
	showTCRConfig()
	toolchain.ShowConfigs()
	language.ShowConfigs()
}

func showTCRConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	utils.Trace("TCR configuration:")
	for _, key := range keys {
		utils.TraceKeyValue(key, viper.Get(key))
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
	Config.VCS = AddVCSParam(cmd)
	Config.MessageSuffix = AddMessageSuffixParam(cmd)
	Config.Trace = AddTraceParam(cmd)
}

// UpdateEngineParams updates TCR engine parameters based on configuration values
func UpdateEngineParams(p *params.Params) {
	p.BaseDir = Config.BaseDir.GetValue()
	p.WorkDir = Config.WorkDir.GetValue()
	p.ConfigDir = Config.ConfigDir.GetValue()
	p.MobTurnDuration = Config.MobTimerDuration.GetValue()
	p.Language = Config.Language.GetValue()
	p.Toolchain = Config.Toolchain.GetValue()
	p.PollingPeriod = Config.PollingPeriod.GetValue()
	p.AutoPush = Config.AutoPush.GetValue()
	p.CommitFailures = Config.CommitFailures.GetValue()
	p.VCS = Config.VCS.GetValue()
	p.MessageSuffix = Config.MessageSuffix.GetValue()
	p.Trace = Config.Trace.GetValue()
}
