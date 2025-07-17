package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"tool/internal/assert"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var environmentFile string
var connStr string
var migrationDir string

var rootCmd = &cobra.Command{
	Use:   "tool",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&environmentFile, "env", "e", ".env", "Environment file location")
	rootCmd.PersistentFlags().StringVarP(&connStr, "conn", "c", "", "Database connection string")
	rootCmd.PersistentFlags().StringVarP(&migrationDir, "migrations", "m", "", "Path to migration directory")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	assert.AssertErr(err)
}

type DatabaseMode string

var (
	READ_ONLY         DatabaseMode = "ro"
	READ_WRITE                     = "rw"
	READ_WRITE_CREATE              = "rwc"
)

type DatabaseOptions struct {
	Mode *DatabaseMode
}

func Database(opts *DatabaseOptions) *sql.DB {
	if len(connStr) == 0 {
		err := godotenv.Load(environmentFile)
		assert.AssertErr(err)
		connStr = os.Getenv("CONN_STR")
	}

	assert.Assert(len(connStr) > 0, "no database connection string")

	if opts != nil {
		var params []string
		if opts.Mode != nil {
			params = append(params, fmt.Sprintf("mode=%s", *opts.Mode))
		}
		connStr = fmt.Sprintf("file:%s?%s", connStr, strings.Join(params, "&"))
	}

	db, err := sql.Open("sqlite3", connStr)
	assert.AssertErr(err)

	return db
}

func MigrationsDir() string {
	if len(migrationDir) == 0 {
		err := godotenv.Load(environmentFile)
		assert.AssertErr(err)
		migrationDir = os.Getenv("MIGRATIONS")
	}

	assert.Assert(len(migrationDir) > 0, "migrations directory was not set")

	_, err := os.Stat(migrationDir)
	assert.Assert(err == nil, "could not find migrations directory")

	return migrationDir
}
