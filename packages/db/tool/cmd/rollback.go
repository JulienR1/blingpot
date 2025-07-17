package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"
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
			if err != nil || count <= 0 {
				fmt.Fprintln(os.Stderr, "Invalid rollback [count] passed in")
				os.Exit(1)
			}
		}

		db := Database(nil)
		defer db.Close()

		if err := migrations.EnsureExists(db); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		stmt, err := db.Prepare("select timestamp, label from migrations order by timestamp desc limit ?;")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not prepare rollback statement:", err)
			os.Exit(1)
		}
		defer stmt.Close()

		rows, err := stmt.Query(count)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not query migrations:", err)
			os.Exit(1)
		}
		defer rows.Close()

		var rollbacks []migrations.Migration
		for rows.Next() {
			var m migrations.Migration
			if err := rows.Scan(&m.Timestamp, &m.Label); err != nil {
				fmt.Fprintln(os.Stderr, "Could not query read migration:", err)
				os.Exit(1)
			}
			rollbacks = append(rollbacks, m)
		}

		tx, err := db.Begin()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not begin rollback transaction:", err)
			os.Exit(1)
		}
		defer tx.Rollback()

		dir := MigrationsDir()
		for i, rollback := range rollbacks {
			filename := fmt.Sprintf("%s-%s.down.sql", rollback.Timestamp, rollback.Label)
			fmt.Fprintf(os.Stderr, "(%d:%d) Executing rollback '%s'\r\n", i+1, len(rollbacks), filename)

			down, err := os.ReadFile(path.Join(dir, filename))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not read rollback file '%s'\r\n", filename)
				os.Exit(1)
			}

			rollback.Down = string(down)
			if err := migrations.Down(tx, rollback); err != nil {
				fmt.Fprintln(os.Stderr, "Could not rollback migration: ", filename)
				os.Exit(1)
			}
		}

		if err := tx.Commit(); err != nil {
			fmt.Fprintln(os.Stderr, "Could not complete rollback")
			os.Exit(1)
		}
	},
}
