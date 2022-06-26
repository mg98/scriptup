package migration

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// NamePattern is the accepted regex pattern for a migration name (excluding timestamp prefix).
const NamePattern = `[a-zA-Z0-9]([a-zA-Z0-9_-]*[a-zA-Z0-9])?`

// DateLayout is the date layout used to timestamp migrations.
const DateLayout = "20060102150405"

// Migration represents a single migration's file and state.
type Migration struct {
	path    string
	Name    string
	Date    time.Time
	UpCmd   *exec.Cmd
	DownCmd *exec.Cmd
}

// New makes a lazy initialisation of a Migration struct.
// Neither does it check if the file exists nor does it parse or validate its contents.
func New(filePath string) (*Migration, error) {
	fileName := filepath.Base(filePath)
	ts, err := time.Parse(DateLayout, fileName[:14])
	if err != nil {
		return nil, err
	}

	return &Migration{
		Name: fileName[:len(fileName)-len(".sh")],
		Date: ts,
		path: filePath,
	}, nil
}

// Parse gets the migration file's content and extracts the Up and Down commands from it
// with respect to the given executor (path to binary).
func (m *Migration) Parse(executor string) error {
	file, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	up, down, err := parse(string(file))
	if err != nil {
		return err
	}
	if up != "" {
		m.UpCmd = exec.Command(executor, "-c", up)
	} else {
		m.UpCmd = nil
	}
	if down != "" {
		m.DownCmd = exec.Command(executor, "-c", down)
	} else {
		m.DownCmd = nil
	}
	return nil
}

func (m *Migration) Up() error {
	return execute(m.UpCmd)
}

func (m *Migration) Down() error {
	return execute(m.DownCmd)
}

func execute(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
