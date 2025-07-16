package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var environmentFile string
var connStr string

var rootCmd = &cobra.Command{
	Use:   "tool",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&environmentFile, "env", "e", "", "Environment file location")
	rootCmd.PersistentFlags().StringVarP(&connStr, "conn", "c", "", "Database connection string")
	rootCmd.MarkFlagsOneRequired("env", "conn")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(checkCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
		if err := godotenv.Load(environmentFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		connStr = os.Getenv("CONN_STR")
	}

	if len(connStr) == 0 {
		fmt.Fprintln(os.Stderr, "no database connection string")
		os.Exit(1)
	}

	if opts != nil {
		var params []string
		if opts.Mode != nil {
			params = append(params, fmt.Sprintf("mode=%s", *opts.Mode))
		}
		connStr = fmt.Sprintf("file:%s?%s", connStr, strings.Join(params, "&"))
	}

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return db
}
