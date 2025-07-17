package migrations

import "fmt"

type Migration struct {
	Timestamp string
	Label     string
	Up        string
	Down      string
}

type MigrationType string

const (
	UP   MigrationType = "up"
	DOWN               = "down"
)

func (m Migration) Title() string {
	return fmt.Sprintf("%s-%s", m.Timestamp, m.Label)
}

func (m Migration) Filename(mode MigrationType) string {
	return fmt.Sprintf("%s.%s.sql", m.Title(), mode)
}
