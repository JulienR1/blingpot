package cmd

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
	"tool/internal/assert"
	"tool/internal/migrations"

	"github.com/spf13/cobra"
)

var contents = []byte("SELECT 1;")

var createCmd = &cobra.Command{
	Use:   "create migration-name",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migration := migrations.Migration{
			Timestamp: strconv.Itoa(int(time.Now().Unix())),
			Label:     args[0],
		}

		dir := MigrationsDir()
		up := migration.Filename(migrations.UP)
		down := migration.Filename(migrations.DOWN)

		err := os.WriteFile(path.Join(dir, up), contents, 0644)
		assert.Assertf(err == nil, "Could not write up migration file: %s\r\n", err)

		err = os.WriteFile(path.Join(dir, down), contents, 0644)
		assert.Assertf(err == nil, "Could not write down migration file: %s\r\n", err)

		fmt.Fprintln(os.Stderr, "Created migration files")
	},
}
