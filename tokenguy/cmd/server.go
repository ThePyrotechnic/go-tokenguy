package cmd

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thepyrotechnic/go-tokenguy/v2/tokenguy"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startServerCmd)

	startServerCmd.Flags().StringP("port", "p", "6666", "Server port")
	startServerCmd.Flags().StringP("host", "a", "0.0.0.0", "Server host address")
	startServerCmd.Flags().StringP("keys-directory", "k", "keys", "Path to directory of public keys to use when validating tokens")
	viper.BindPFlag("server.host", startServerCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", startServerCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.keys", startServerCmd.Flags().Lookup("keys-directory"))
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the tokenguy server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Error: Must specify a sub-command")
		fmt.Println()
		cmd.Help()
	},
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the tokenguy server",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("server.host")
		port := viper.GetString("server.port")

		keysDir := viper.GetString("server.keys")
		fileinfo, err := os.Stat(keysDir)
		if err != nil {
			panic(err)
		}
		if !fileinfo.IsDir() {
			panic(fmt.Errorf("Provided keys directory is not a directory"))
		}
		matches, err := filepath.Glob(filepath.Join(keysDir, "*"))
		if err != nil {
			panic(err)
		}
		if matches == nil {
			panic(fmt.Errorf("Provided keys directory is empty"))
		}

		keyMap := make(map[string]*rsa.PublicKey)
		for a := 0; a < len(matches); a++ {
			data, err := os.ReadFile(matches[a])
			if err != nil {
				log.Println(matches[a], ": ", err)
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM(data)
			if err != nil {
				log.Println(matches[a], ": ", err)
			}
			keyMap[filepath.Base(matches[a])] = key
		}

		router := tokenguy.Router(keyMap)
		fmt.Printf("Starting API server on http://%s:%s ...\n", host, port)
		if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
			panic(fmt.Errorf("Error starting server: %s", err))
		}
	},
}
