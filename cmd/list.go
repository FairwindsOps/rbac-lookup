package cmd

import (
	"github.com/reactiveops/rbac-lookup/lookup"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all RBAC Bindings by user",
	Run: func(cmd *cobra.Command, args []string) {
		lookup.ListAll()
	},
}
