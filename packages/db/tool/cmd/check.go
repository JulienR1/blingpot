package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate connection to database",
	Run: func(cmd *cobra.Command, args []string) {
		db := Database(&DatabaseOptions{Mode: &READ_ONLY})
		defer db.Close()

		if db.Ping() == nil {
			fmt.Fprintln(os.Stderr, "Connection to database was successful")
		} else {
			fmt.Fprintln(os.Stderr, "Could not connect to database")
		}
	},
}
