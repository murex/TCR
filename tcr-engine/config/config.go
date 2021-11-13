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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// SaveConfigFlag is a placeholder for command line parameter allowing to save configuration to a file
var SaveConfigFlag bool

// Init initializes viper configuration
func Init(configFile string) {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, _ := homedir.Dir()
		//cobra.CheckErr(err)

		// Search config in home directory with name "tcr.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("tcr.yaml")
		viper.SetConfigFile(filepath.Join(home, "tcr.yaml"))
	}

	viper.SetEnvPrefix(ApplicationName)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		trace("Loading configuration: ", viper.ConfigFileUsed())
	}
}

// Save saves viper configuration
func Save() {
	if SaveConfigFlag {
		trace("Saving configuration: ", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			if os.IsNotExist(err) {
				_ = viper.WriteConfigAs(viper.ConfigFileUsed())
			} else {
				trace("Error while saving configuration file: ", err)
			}
		}
	}
}

func trace(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, "["+ApplicationName+"]", fmt.Sprint(a...))
}
