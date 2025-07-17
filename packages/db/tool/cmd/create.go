package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var contents = []byte("SELECT 1;")

var createCmd = &cobra.Command{
	Use:   "create migration-name",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		timestamp := strconv.Itoa(int(time.Now().Unix()))

		dir := MigrationsDir()
		up := fmt.Sprintf("%s-%s.up.sql", timestamp, migrationName)
		down := fmt.Sprintf("%s-%s.down.sql", timestamp, migrationName)

		if err := os.WriteFile(path.Join(dir, up), contents, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Could not write up migration file", err)
			os.Exit(1)
		}

		if err := os.WriteFile(path.Join(dir, down), contents, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "Could not write down migration file", err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, "Created migration files")
	},
}
