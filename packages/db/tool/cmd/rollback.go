package cmd

import (
	"log"
	"os"
	"path"
	"strconv"
	"tool/internal/assert"
	"tool/internal/migrations"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [count]",
	Short: "Execute [count] down migration files",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var count = 1

		if len(args) > 0 {
			count, err = strconv.Atoi(args[0])
			assert.Assert(err == nil && count > 0, "Invalid rollback [count] passed in")
		}

		db := Database(nil)
		defer db.Close()

		assert.AssertErr(migrations.EnsureExists(db))

		stmt, err := db.Prepare("select timestamp, label from migrations order by timestamp desc limit ?;")
		assert.Assertf(err == nil, "Could not prepare rollback statement: %s\r\n", err)
		defer stmt.Close()

		rows, err := stmt.Query(count)
		assert.Assertf(err == nil, "Could not query migrations: %s\r\n", err)
		defer rows.Close()

		var rollbacks []migrations.Migration
		for rows.Next() {
			var m migrations.Migration
			err := rows.Scan(&m.Timestamp, &m.Label)
			assert.Assertf(err == nil, "Could not query read migration: %s\r\n", err)
			rollbacks = append(rollbacks, m)
		}

		tx, err := db.Begin()
		assert.Assertf(err == nil, "Could not begin rollback transaction: %s\r\n", err)
		defer tx.Rollback()

		dir := MigrationsDir()
		for i, rollback := range rollbacks {
			filename := rollback.Filename(migrations.DOWN)
			log.Printf("(%d:%d) Executing rollback '%s'\r\n", i+1, len(rollbacks), filename)

			down, err := os.ReadFile(path.Join(dir, filename))
			assert.Assertf(err == nil, "Could not read rollback file '%s'\r\n", filename)

			rollback.Down = string(down)
			err = migrations.Down(tx, rollback)
			assert.Assertf(err == nil, "Could not rollback migration: %s\r\n", filename)
		}

		err = tx.Commit()
		assert.Assert(err == nil, "Could not complete rollback")
	},
}
