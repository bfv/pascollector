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
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/bfv/pascollector/misc"
	"github.com/bfv/pascollector/types"
)

var forceOverwrite bool = false

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("config called")
	// },
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup a new configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		setup()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "shows the configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		showConfiguration()
	},
}

func init() {

	configCmd.AddCommand(setupCmd)
	configCmd.AddCommand(showCmd)

	rootCmd.AddCommand(configCmd)

}

func getDefaultConfiguration() types.ConfigFile {

	config := types.ConfigFile{}
	config.ClientId = "<client-id>"
	config.Server = "<server-name>"
	config.Tag = ""
	config.Port = DefaultPort
	config.CollectInterval = 60
	config.SendInterval = 60
	config.PasInstances = []types.PasInstance{
		{
			Name: "oepas1",
			Url:  "http://localhost:8810",
		},
	}

	return config
}

func showConfiguration() {
	config, _ := yaml.Marshal(Config)
	fmt.Println(string(config))
}

func setup() {

	// check user first, on Linux it should be root
	misc.CheckUser()

	programConfigDir := misc.GetConfigDir()
	configFilename := misc.GetConfigurationFilename()

	_, err := os.Stat(programConfigDir)
	if os.IsNotExist(err) {
		os.Mkdir(programConfigDir, 0777)
	}

	if _, err = os.Stat(configFilename); err == nil && !forceOverwrite {
		log.Fatal(configFilename + " exists, use -f to overwrite")
	}

	config, _ := yaml.Marshal(getDefaultConfiguration())
	err = ioutil.WriteFile(configFilename, config, 0744)
	if err != nil {
		log.Fatalln("unable to create .pascollector.yaml in " + programConfigDir)
	}

	fmt.Println(configFilename + " created")
}
