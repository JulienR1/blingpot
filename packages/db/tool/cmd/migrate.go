package cmd

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tool/internal/assert"
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
			assert.Assert(err == nil && count > 0, "Invalid migration [count] passed in")
		}

		db := Database(nil)
		defer db.Close()

		assert.AssertErr(migrations.EnsureExists(db))

		var latestMigration migrations.Migration
		_ = db.
			QueryRow("select timestamp, label from migrations order by timestamp desc limit 1;").
			Scan(&latestMigration.Timestamp, &latestMigration.Label)
		var latestMigrationFileName = latestMigration.Filename(migrations.UP)

		var migrationsToExecute []migrations.Migration
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
			log.Println("Already up to date, no migrations to execute.")
			return
		}

		tx, err := db.Begin()
		assert.Assertf(err == nil, "Could not begin migration transaction: %s\r\n", err)
		defer tx.Rollback()

		for i, migration := range migrationsToExecute {
			log.Printf("(%d:%d) Executing up migration '%s'\r\n", i+1, len(migrationsToExecute), migration.Filename(migrations.UP))
			err := migrations.Up(tx, migration)
			assert.AssertErr(err)
		}

		err = tx.Commit()
		assert.Assertf(err == nil, "Could not complete migration: %s\r\n", err)
	},
}
