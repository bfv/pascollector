/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/bfv/pascollector/misc"
	"github.com/spf13/cobra"
)

// passwdCmd represents the passwd command
var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "A brief description of your command",
	Long:  ``,
}

// passwdCmd represents the passwd command
var passwdEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		encr := misc.Encrypt(args[0])
		fmt.Println(encr)
		// fmt.Println("decrypted: " + misc.Decrypt(encr))
	},
}

func init() {
	passwdCmd.AddCommand(passwdEncryptCmd)
	rootCmd.AddCommand(passwdCmd)
}
