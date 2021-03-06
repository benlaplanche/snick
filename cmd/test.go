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
	"time"

	"github.com/benlaplanche/snick/environment"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
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

type appEnv struct {
	timestamp   time.Time
	debug       bool
	rules       string
	output      string
	files       string
	results     []Result
	environment string
}

func (app *appEnv) setup(cmd *cobra.Command) error {
	app.files, _ = cmd.Flags().GetString("files")
	app.rules, _ = cmd.Flags().GetString("rego")
	app.debug, _ = cmd.Flags().GetBool("debug")
	app.environment = environment.DetectENV()
	return nil
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "evaluate config files against policies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var app appEnv

		err := app.setup(cmd)
		if err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadFile(app.files)
		if err != nil {
			log.Fatal(err)
		}

		if app.debug {
			fmt.Println("file path:", app.files)
			fmt.Println("rego path:", app.rules)
			fmt.Println("file contents:", string(content))

			json, err := yaml.YAMLToJSON(content)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("json:", string(json))
		}

		fmt.Printf("Environment detected as: %s \n", app.environment)

		err = yaml.Unmarshal(content, &j)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		r := rego.New(
			rego.Query("data.main"),
			rego.Load([]string{app.rules}, nil))

		query, err := r.PrepareForEval(ctx)
		if err != nil {
			log.Fatal(err)
		}

		rs, err := query.Eval(ctx, rego.EvalInput(j))
		if err != nil {
			log.Fatal(err)
		}

		output := prepareOutput(rs)
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
