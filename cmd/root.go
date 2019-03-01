// Copyright 2018 ReactiveOps
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

	"github.com/reactiveops/rbac-lookup/lookup"
	"github.com/spf13/cobra"
)

var outputFormat string
var enableGke bool
var kubeContext string
var subjectType string

var rootCmd = &cobra.Command{
	Use:   "rbac-lookup [subject query]",
	Short: "rbac-lookup provides a simple way to view RBAC bindings by user",
	Long:  "rbac-lookup provides a missing Kubernetes API to view RBAC bindings by user",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.ParseFlags(args); err != nil {
			fmt.Printf("Error parsing flags: %v", err)
		}

		lookup.List(args, kubeContext, outputFormat, subjectType, enableGke)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format (normal,wide)")
	rootCmd.PersistentFlags().StringVarP(&kubeContext, "context", "", "", "context to use for Kubernetes config")
	rootCmd.PersistentFlags().StringVarP(&subjectType, "type", "t", "", "filter by this RBAC subject type")
	rootCmd.PersistentFlags().BoolVar(&enableGke, "gke", false, "enable GKE integration")
}

// Execute is the primary entrypoint for this CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
