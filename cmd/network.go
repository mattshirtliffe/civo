// Copyright © 2016 Absolute DevOps Ltd <info@absolutedevops.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/absolutedevops/civo/api"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var networkFullIDs bool

// networkCmd represents the accounts command
var networkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"networks"},
	Short:   "List all networks",
	Long:    `List the networks for the current account`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.NetworksList()
		if err != nil {
			errorColor := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Println(errorColor("An error occured:"), err.Error())
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetHeader([]string{"ID", "Label", "Default"})
		items, _ := result.Children()
		for _, child := range items {
			var id string
			if networkFullIDs {
				id = child.S("id").Data().(string)
			} else {
				parts := strings.Split(child.S("id").Data().(string), "-")
				id = parts[0]
			}

			defaultLabel := "no"
			if child.S("default").Data().(bool) {
				defaultLabel = "yes"
			}

			table.Append([]string{
				id,
				child.S("label").Data().(string),
				defaultLabel,
			})
		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(networkCmd)
	networkCmd.Flags().BoolVarP(&networkFullIDs, "full-ids", "f", false, "Return full IDs for networks")
}
