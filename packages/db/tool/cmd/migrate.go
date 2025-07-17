package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tool/internal/migrations"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [count]",
	Short: "Execute migration files or up to [count] files",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var count = 0

		if len(args) > 0 {
			count, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid migration [count] passed in")
				os.Exit(1)
			}
		}

		db := Database(nil)
		defer db.Close()

		if err := migrations.EnsureExists(db); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var migrationsToExecute []migrations.Migration

		var latestMigration migrations.Migration
		_ = db.
			QueryRow("select timestamp, label from migrations order by timestamp desc limit 1;").
			Scan(&latestMigration.Timestamp, &latestMigration.Label)
		var latestMigrationFileName = fmt.Sprintf("%s-%s.up.sql", latestMigration.Timestamp, latestMigration.Label)

		var walkingNewMigrations = len(latestMigration.Timestamp) == 0
		filepath.WalkDir(MigrationsDir(), func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() || strings.HasSuffix(d.Name(), ".up.sql") == false {
				return nil
			}

			if walkingNewMigrations {
				up, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				down, err := os.ReadFile(strings.Replace(path, ".up.sql", ".down.sql", 1))
				if err != nil {
					return err
				}

				parts := strings.SplitAfterN(d.Name(), "-", 2)
				if len(parts) != 2 {
					return errors.New("filename does not match intended format")
				}

				m := migrations.Migration{
					Timestamp: parts[0][:len(parts[0])-1],
					Label:     strings.Replace(parts[1], ".up.sql", "", 1),
					Up:        string(up),
					Down:      string(down),
				}
				migrationsToExecute = append(migrationsToExecute, m)
			}

			if count > 0 && len(migrationsToExecute) == count {
				walkingNewMigrations = false
				return fs.SkipAll
			}

			walkingNewMigrations = walkingNewMigrations || latestMigrationFileName == d.Name()
			return nil
		})

		if len(migrationsToExecute) == 0 {
			fmt.Fprintf(os.Stderr, "Already up to date, no migrations to execute.\r\n")
			return
		}

		tx, err := db.Begin()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not begin migration transaction:", err)
			os.Exit(1)
		}
		defer tx.Rollback()

		for i, migration := range migrationsToExecute {
			fmt.Fprintf(os.Stderr, "(%d:%d) Executing up migration '%s'\r\n", i+1, len(migrationsToExecute), migration.Label)
			if err := migrations.Up(tx, migration); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		if err := tx.Commit(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not complete migration", err)
			os.Exit(1)
		}
	},
}
