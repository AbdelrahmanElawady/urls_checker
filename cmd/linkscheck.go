/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rawdaGastan/urls_checker/pkg"
)

// linkscheckCmd represents the linkscheck command
var linkscheckCmd = &cobra.Command{
	Use:   "linkscheck",
	Short: "Checks the status of a website's urls",
	Long:  `Iterates through the urls of a given website and checks the status of each url`,
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := cmd.Flags().GetString("config")
		configFile := strings.Split(config, ".")

		if len(configFile) > 2 {
			fmt.Println("error: Invalid file name")
			return
		}

		if configFile[1] != "toml" {
			fmt.Println("error: Invalid file extension. should be .toml")
			return
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("error: ", err)
			return
		}

		// parse toml file
		viper.AddConfigPath(currentDir)
		viper.SetConfigName(configFile[0])
		viper.SetConfigType(configFile[1])

		errParseFile := viper.ReadInConfig()

		if errParseFile != nil {
			fmt.Println("error: ", errParseFile)
			return
		}

		var sites map[string]interface{} = viper.GetStringMap("sites")

		if len(sites) == 0 {
			fmt.Println("error: no sites provided")
			return
		}

		for _, site := range sites {
			parsed := site.(map[string]interface{})
			url := parsed["url"].(string)

			if url == "" {
				fmt.Println("error: no url provided")
				return
			}
			checkErr := pkg.Check(url)

			if checkErr != nil {
				fmt.Println("error: ", checkErr)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(linkscheckCmd)
}
