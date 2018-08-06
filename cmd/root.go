package cmd

import (
	"fmt"
	"os"

	"github.com/reactiveops/rbac-lookup/lookup"
	"github.com/spf13/cobra"
)

var output string

var rootCmd = &cobra.Command{
	Use:   "rbac-lookup [subject query]",
	Short: "rbac-lookup provides a simple way to view RBAC bindings by user",
	Long:  "rbac-lookup provides a missing Kubernetes API to view RBAC bindings by user",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.ParseFlags(args); err != nil {
			fmt.Printf("Error parsing flags: %v", err)
		}

		lookup.List(args, output)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format (normal,wide)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
