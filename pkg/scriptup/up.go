package scriptup

import (
	"fmt"
	"github.com/mg98/scriptup/pkg/scriptup/migration"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
)

func MigrateUp(cfg *Config, steps int) error {
	s, err := cfg.InitMigrationState()
	if err != nil {
		return err
	}
	defer s.Close()

	latestMigration, err := s.Latest()
	if err != nil {
		return err
	}

	var migrations []*migration.Migration
	if err := filepath.WalkDir(cfg.Directory, func(s string, f fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !f.IsDir() && regexp.MustCompile(
			fmt.Sprintf("^[0-9]{14}_%s\\.sh$", migration.NamePattern),
		).MatchString(f.Name()) {
			m, err := migration.New(cfg.Directory + "/" + f.Name())
			if err != nil {
				return err
			}
			if latestMigration == nil || m.Date.Unix() > latestMigration.Date.Unix() {
				migrations = append(migrations, m)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// sort by time ascending
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Date.Before(migrations[j].Date)
	})

	if steps >= 0 {
		migrations = migrations[:steps]
	}

	for _, m := range migrations {
		// DISCUSS: Parse happens in MigrateDown too. Should this be part of NewMigrator or even an extra func?
		if err := m.Parse(cfg.Executor); err != nil {
			return err
		}
		if err := m.Up(); err != nil {
			return err
		}
		if err = s.Add(m); err != nil {
			return err
		}
		fmt.Printf("[+] %s\n", m.Name)
	}

	return nil
}
