package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration file",
	Run: func(cmd *cobra.Command, args []string) {
		Database(nil)
	},
}
