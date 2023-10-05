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
	"github.com/thepyrotechnic/go-tokenguy/v2/tokenguy"
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("keys-directory", "k", "keys", "Path to directory of public keys to use when validating tokens")
	viper.BindPFlag("server.keys", validateCmd.Flags().Lookup("keys-directory"))
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate a JWT passed to stdin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
		if tokenguy.Validate(tokenguy.GetKeys(), args[0]) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	},
}
