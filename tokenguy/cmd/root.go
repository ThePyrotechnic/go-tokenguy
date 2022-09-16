/*
tokenguy is a web server which validates JWTs
Copyright (C) 2022  Michael Manis

  This program is free software; you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation; either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program; if not, write to the Free Software Foundation,
  Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd *cobra.Command
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tokenguy",
		Short: "JWT token helper",
		Long:  `TODO`,
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd = NewRootCmd()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/tokenguy/config)")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("TOKENGUY")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")

		if homedir := os.Getenv("HOME"); homedir == "" {
			panic(fmt.Errorf("$HOME is not set. Either set $HOME or manually specify a config directory"))
		}

		viper.AddConfigPath("$HOME/.config/tokenguy/")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Error parsing config file: %w", err))
		}
	}
}
