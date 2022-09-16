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
