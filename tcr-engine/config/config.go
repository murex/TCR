/*
Copyright (c) 2021 Murex

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
	"github.com/murex/tcr/tcr-engine/engine"
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sort"
)

// TcrConfig wraps all possible TCR configuration parameters
type TcrConfig struct {
	BaseDir          *StringParam
	ConfigDir        *StringParam
	Language         *StringParam
	Toolchain        *StringParam
	PollingPeriod    *DurationParam
	MobTimerDuration *DurationParam
	AutoPush         *BoolParam
}

func (c TcrConfig) reset() {
	c.Language.reset()
	c.Toolchain.reset()
	c.PollingPeriod.reset()
	c.MobTimerDuration.reset()
	c.AutoPush.reset()
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

// Init initializes TCR configuration
func Init() {

	// Initialize configuration directory path
	initConfigDirPath()

	// Make sure configuration directory exists
	createConfigDir()

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
	}
}

func initConfigDirPath() {
	if Config.ConfigDir.GetValue() == "" {
		// If configuration directory is not specified, we use by default the current directory
		configDirPath = filepath.Join(".", configDirRoot)
	} else {
		configDirPath = filepath.Join(Config.ConfigDir.GetValue(), configDirRoot)
	}
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
	Config.reset()
	Save()
}

// Show displays current TCR configuration
func Show() {
	showConfigOrigin()
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		trace("- ", key, ": ", viper.Get(key))
	}
}

func showConfigOrigin() {
	_, err := os.Stat(viper.ConfigFileUsed())
	if os.IsNotExist(err) {
		trace("No configuration file found. Showing default configuration")
	} else {
		trace("Configuration file: ", viper.ConfigFileUsed())
	}
}

func trace(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, "["+settings.ApplicationName+"]", fmt.Sprint(a...))
}

// AddParameters add parameter to the provided command cmd
func AddParameters(cmd *cobra.Command, defaultBaseDir string) {
	Config.BaseDir = AddBaseDirParamWithDefault(cmd, defaultBaseDir)
	Config.ConfigDir = AddConfigDirParam(cmd)
	Config.Language = AddLanguageParam(cmd)
	Config.Toolchain = AddToolchainParam(cmd)
	Config.PollingPeriod = AddPollingPeriodParam(cmd)
	Config.MobTimerDuration = AddMobTimerDurationParam(cmd)
	Config.AutoPush = AddAutoPushParam(cmd)
}

// UpdateEngineParams updates TCR engine parameters based on configuration values
func UpdateEngineParams(params *engine.Params) {
	params.BaseDir = Config.BaseDir.GetValue()
	params.ConfigDir = Config.ConfigDir.GetValue()
	params.MobTurnDuration = Config.MobTimerDuration.GetValue()
	params.Language = Config.Language.GetValue()
	params.Toolchain = Config.Toolchain.GetValue()
	params.PollingPeriod = Config.PollingPeriod.GetValue()
	params.AutoPush = Config.AutoPush.GetValue()
}
