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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/benlaplanche/snick/environment"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/jmespath/go-jmespath"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var files string
var regoFiles string
var debug bool
var j interface{}
var output []Result
var filter string

type Result struct {
	id       string `json:"id"`
	name     string `json:"name"`
	severity string `json:"severity"`
	response string `json:"response"`
	status   string `json:"status"`
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "evaluate config files against policies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// setup
		filesPath, _ := cmd.Flags().GetString("files")
		regoPath, _ := cmd.Flags().GetString("rego")
		debug, _ := cmd.Flags().GetBool("debug")
		filter, _ := cmd.Flags().GetString("filter")

		content, err := ioutil.ReadFile(filesPath)
		if err != nil {
			log.Fatal(err)
		}

		if debug {
			fmt.Println("file path:", filesPath)
			fmt.Println("rego path:", regoPath)
			fmt.Println("file contents:", string(content))
			fmt.Println("filter:", filter)

			json, err := yaml.YAMLToJSON(content)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("json:", string(json))
		}

		fmt.Printf("Environment detected as: %s \n", environment.DetectENV())

		err = yaml.Unmarshal(content, &j)
		if err != nil {
			log.Fatal(err)
		}

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

		output := prepareOutput(rs)

		if filter != "" {
			data, err := json.Marshal(&output)
			fmt.Println(string(data))

			if err != nil {
				fmt.Println("error marshaling to json")
			}

			blah, err := jmespath.Search(filter, data)
			fmt.Println(string(data))
			fmt.Println(blah)
			if err != nil {
				fmt.Println("error using th filter")
			}
		}
		tableOutput(output)

	},
}

func prepareOutput(rs rego.ResultSet) []Result {

	for _, result := range rs {
		for _, expression := range result.Expressions {

			var expressionValues map[string]interface{}
			if _, ok := expression.Value.(map[string]interface{}); ok {
				expressionValues = expression.Value.(map[string]interface{})
			}

			if len(expressionValues) == 0 {
				output = append(output, Result{})
				continue
			}

			for _, v := range expressionValues {

				var nestedValues []interface{}
				if _, ok := v.([]interface{}); ok {
					nestedValues = v.([]interface{})
				}

				for _, x := range nestedValues {

					switch val := x.(type) {
					// Policies that only return a single string (e.g. deny[msg])
					case string:
						// result := output.Result{
						// 	Message: val,
						// }
						// output = append(output, result)

					// Policies that return metadata (e.g. deny[{"msg": msg}])
					case map[string]interface{}:
						var rx string
						if val["status"].(string) == "allow" {
							rx = val["allow_response"].(string)
						} else {
							rx = val["deny_response"].(string)
						}

						result := Result{
							id:       val["id"].(string),
							status:   val["status"].(string),
							name:     val["name"].(string),
							response: rx,
							severity: val["severity"].(string),
						}
						output = append(output, result)
					}
				}
			}
		}
	}
	return output
}

func tableOutput(output []Result) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Status", "ID", "Name", "Severity", "Response")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, r := range output {
		tbl.AddRow(r.status, r.id, r.name, r.severity, r.response)
	}

	tbl.Print()
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&files, "files", "f", "", "location to config files to be tested")
	testCmd.MarkPersistentFlagRequired("files")

	testCmd.PersistentFlags().StringVarP(&regoFiles, "rego", "r", "", "location to rego files")
	testCmd.MarkPersistentFlagRequired("rego")

	testCmd.PersistentFlags().StringVarP(&filter, "filter", "", "", "filter the results")

	testCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
}
