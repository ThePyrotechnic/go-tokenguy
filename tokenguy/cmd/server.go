/*
tokenguy is a web server which validates JWTs
Copyright (C) 2024  Michael Manis

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
	"github.com/thepyrotechnic/go-tokenguy/v2/tokenguy"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startServerCmd)

	startServerCmd.Flags().StringP("port", "p", "6666", "Server port")
	startServerCmd.Flags().StringP("host", "a", "0.0.0.0", "Server host address")
	startServerCmd.Flags().StringP("public-keys", "k", "keys/public", "Path to directory of public keys to use when validating tokens")
	startServerCmd.Flags().StringP("private-keys", "e", "keys/private", "Path to directory of private keys to use when signing tokens")
	viper.BindPFlag("server.host", startServerCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", startServerCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.keys.public", startServerCmd.Flags().Lookup("public-keys"))
	viper.BindPFlag("server.keys.private", startServerCmd.Flags().Lookup("private-keys"))
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

		pubKeys := tokenguy.GetPublicKeys()
		privKeys := tokenguy.GetPrivateKeys()

		router := tokenguy.Router(pubKeys, privKeys)
		fmt.Printf("Starting API server on http://%s:%s ...\n", host, port)
		if err := router.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
			panic(fmt.Errorf("error starting server: %s", err))
		}
	},
}
