package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rbac-lookup [subject/query]",
	Short: "rbac-lookup provides a simple way to view RBAC bindings by user",
	Long:  "rbac-lookup provides a missing Kubernetes API to view RBAC bindings by user",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world " + strings.Join(args, " "))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
