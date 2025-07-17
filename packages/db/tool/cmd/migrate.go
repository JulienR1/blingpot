package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"tool/internal/migrations"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [count]",
	Short: "Execute migration files or up to [count] files",
	Run: func(cmd *cobra.Command, args []string) {
		db := Database(nil)
		defer db.Close()

		if err := migrations.EnsureExists(db); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// diff fs w/ db
		migrationsToExecute := []migrations.Migration{}

		filepath.WalkDir(MigrationsDir(), func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() == false {
				// TODO: check if to execute
			}
			return nil
		})

		var migrationIndex = 0
		var mustRevertMigrations = false

		for ; migrationIndex < len(migrationsToExecute); migrationIndex++ {
			migration := migrationsToExecute[migrationIndex]
			fmt.Fprintf(os.Stderr, "(%d:%d) Executing up migration '%s'\r\n", migrationIndex, len(migrationsToExecute), migration.Label)

			if err := migrations.Up(db, migration); err != nil {
				fmt.Fprintln(os.Stderr, err)
				mustRevertMigrations = true
				break
			}
		}

		if mustRevertMigrations {
			for _, migration := range slices.Backward(migrationsToExecute[:migrationIndex]) {
				if err := migrations.Down(db, migration); err != nil {
					fmt.Fprintln(os.Stderr, "Could not restore database to previous state", err)
					os.Exit(1)
				}
			}
		}
	},
}
