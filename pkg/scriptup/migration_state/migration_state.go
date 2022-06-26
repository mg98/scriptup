package migration_state

import (
	"fmt"
	"github.com/mg98/scriptup/pkg/scriptup/migration"
	"github.com/mg98/scriptup/pkg/scriptup/storage"
	"io/fs"
	"path/filepath"
	"regexp"
)

type MigrationState struct {
	migrationDir string
	storage      storage.Storage
}

func New(migrationDir string, storage storage.Storage) *MigrationState {
	return &MigrationState{migrationDir: migrationDir, storage: storage}
}

func (s *MigrationState) Add(m *migration.Migration) error {
	return s.storage.Append(m.Name)
}

func (s *MigrationState) DropLatest() error {
	return s.storage.Pop()
}

func (s *MigrationState) All(o storage.Order) ([]*migration.Migration, error) {
	entries, err := s.storage.All(o)
	if err != nil {
		return nil, err
	}
	var ms []*migration.Migration
	for _, entry := range entries {
		m, err := migration.New(s.migrationPath(entry))
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// Files returns the list of migration files identified on the disk in ascending order.
func (s *MigrationState) Files() ([]*migration.Migration, error) {
	var migrations []*migration.Migration
	if err := filepath.WalkDir(s.migrationDir, func(_ string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && regexp.MustCompile(
			fmt.Sprintf("^[0-9]{14}_%s\\.sh$", migration.NamePattern),
		).MatchString(f.Name()) {
			m, err := migration.New(s.migrationDir + "/" + f.Name())
			if err != nil {
				return err
			}
			migrations = append(migrations, m)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return migrations, nil
}

func (s *MigrationState) Latest() (*migration.Migration, error) {
	name, err := s.storage.Latest()
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, nil
	}
	return migration.New(s.migrationPath(*name))
}

func (s *MigrationState) migrationPath(name string) string {
	return fmt.Sprintf("%s/%s.sh", s.migrationDir, name)
}

func (s *MigrationState) Close() error {
	return s.storage.Close()
}
