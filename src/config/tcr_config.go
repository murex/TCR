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
	"github.com/murex/tcr/helpers"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/settings"
	"github.com/murex/tcr/toolchain"
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
	GitRemote        *StringParam
	AutoPush         *BoolParam
	Variant          *StringParam
	VCS              *StringParam
	MessageSuffix    *StringParam
	Trace            *StringParam
	PortNumber       *IntParam
}

func (c TcrConfig) reset() {
	c.Language.reset()
	c.Toolchain.reset()
	c.PollingPeriod.reset()
	c.MobTimerDuration.reset()
	c.AutoPush.reset()
	c.GitRemote.reset()
	c.Variant.reset()
	c.VCS.reset()
	c.MessageSuffix.reset()
	c.Trace.reset()
	c.PortNumber.reset()
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
	helpers.SetSimpleTrace(writer)
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
		helpers.Trace("Loading configuration: ", viper.ConfigFileUsed())
	} else {
		helpers.Trace("No configuration file found")
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
		helpers.Trace("Creating TCR configuration directory: ", configDirPath)
		err := os.MkdirAll(configDirPath, 0750)
		if err != nil {
			helpers.Trace("Error creating TCR configuration directory: ", err)
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
	helpers.Trace("Saving configuration: ", viper.ConfigFileUsed())
	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			_ = viper.WriteConfigAs(viper.ConfigFileUsed())
		} else {
			helpers.Trace("Error while saving configuration file: ", err)
		}
	}
}

// Reset resets TCR configuration to default value
func Reset() {
	helpers.Trace("Resetting configuration to default values")
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
	helpers.Trace()
	showTCRConfig()
	toolchain.ShowConfigs()
	language.ShowConfigs()
}

func showTCRConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	helpers.Trace("TCR configuration:")
	for _, key := range keys {
		helpers.TraceKeyValue(key, viper.Get(key))
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
	Config.GitRemote = AddGitRemoteParam(cmd)
	Config.AutoPush = AddAutoPushParam(cmd)
	Config.Variant = AddVariantParam(cmd)
	Config.VCS = AddVCSParam(cmd)
	Config.MessageSuffix = AddMessageSuffixParam(cmd)
	Config.Trace = AddTraceParam(cmd)
	Config.PortNumber = AddPortNumberParam(cmd)
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
	p.GitRemote = Config.GitRemote.GetValue()
	p.Variant = Config.Variant.GetValue()
	p.VCS = Config.VCS.GetValue()
	p.MessageSuffix = Config.MessageSuffix.GetValue()
	p.Trace = Config.Trace.GetValue()
	p.PortNumber = Config.PortNumber.GetValue()
}
