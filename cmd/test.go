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
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
)

var files string
var regoFiles string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "evaluate config files against policies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")

		filesPath, _ := cmd.Flags().GetString("files")
		if filesPath != "" {
			fmt.Println("file path is ", filesPath)
		}

		regoPath, _ := cmd.Flags().GetString("rego")
		if regoPath != "" {
			fmt.Println("rego path is ", regoPath)
		}

		content, err := ioutil.ReadFile(filesPath)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(content))

		json, err := yaml.YAMLToJSON(content)
		if err != nil {
			log.Fatal(err)
		}

		var j interface{}
		err = yaml.Unmarshal(content, &j)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(json)
		fmt.Println(string(json))

		ctx := context.Background()

		r := rego.New(
			rego.Query("data.main"),
			rego.Load([]string{regoPath}, nil))

		query, err := r.PrepareForEval(ctx)
		if err != nil {
			log.Fatal(err)
		}

		rs, err := query.Eval(ctx, rego.EvalInput(j))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(rs)

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&files, "files", "f", "", "location to config files to be tested")
	testCmd.MarkPersistentFlagRequired("files")

	testCmd.PersistentFlags().StringVarP(&regoFiles, "rego", "r", "", "location to rego files")
	testCmd.MarkPersistentFlagRequired("rego")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
