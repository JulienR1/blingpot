package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate connection to database",
	Run: func(cmd *cobra.Command, args []string) {
		db := Database(&DatabaseOptions{Mode: &READ_ONLY})
		defer db.Close()

		if db.Ping() == nil {
			log.Println("Connection to database was successful")
		} else {
			log.Println("Could not connect to database")
		}
	},
}
