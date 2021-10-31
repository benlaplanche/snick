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
var debug bool
var j interface{}

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

		content, err := ioutil.ReadFile(filesPath)
		if err != nil {
			log.Fatal(err)
		}

		if debug {
			fmt.Println("file path:", filesPath)
			fmt.Println("rego path:", regoPath)
			fmt.Println("file contents:", string(content))

			json, err := yaml.YAMLToJSON(content)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("json:", string(json))
		}

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

		prettyOutput(rs)

	},
}

func prettyOutput(rs rego.ResultSet) {
	// top := rs[0].Expressions
	// fmt.Println(reflect.TypeOf(top).Kind())
	// val := top[0].Value
	// fmt.Println(reflect.TypeOf(val).Kind())

	for _, result := range rs {
		// fmt.Println(result)
		for _, expression := range result.Expressions {
			fmt.Println(expression.Value)
			// fmt.Println(expression.Value.(map[string]interface{}))
			fmt.Println("\n")
			// Rego rules that are intended for evaluation should return a slice of values.
			// For example, deny[msg] or violation[{"msg": msg}].
			//
			// When an expression does not have a slice of values, the expression did not
			// evaluate to true, and no message was returned.
			var expressionValues map[string]interface{}
			if _, ok := expression.Value.(map[string]interface{}); ok {
				expressionValues = expression.Value.(map[string]interface{})
			}
			// fmt.Println(expressionValues...)
			if len(expressionValues) == 0 {
				// results = append(results, output.Result{})
				// fmt.Println("blah")
				continue
			}

			for _, v := range expressionValues {
				fmt.Println(v)
				switch val := v.(type) {

				// Policies that only return a single string (e.g. deny[msg])
				case string:
					// result := output.Result{
					// 	Message: val,
					// }
					// results = append(results, result)
					fmt.Println(v)
					fmt.Println(val)

				// Policies that return metadata (e.g. deny[{"msg": msg}])
				case map[string]interface{}:
					// result, err := output.NewResult(val)
					// if err != nil {
					// return output.QueryResult{}, fmt.Errorf("new result: %w", err)
					// fmt.Println("error")
					// }

					// results = append(results, result)
					fmt.Println(v)
				}
			}
		}
	}

	// headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	// columnFmt := color.New(color.FgYellow).SprintfFunc()

	// tbl := table.New("Status", "ID", "Name", "Severity", "Response")
	// tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// for _, r := range allow {
	// 	tbl.AddRow(r.status, r.id, r.name, r.severity, r.allow_response)
	// }

	// tbl.Print()
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&files, "files", "f", "", "location to config files to be tested")
	testCmd.MarkPersistentFlagRequired("files")

	testCmd.PersistentFlags().StringVarP(&regoFiles, "rego", "r", "", "location to rego files")
	testCmd.MarkPersistentFlagRequired("rego")

	testCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
