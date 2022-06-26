package scriptup

import (
	"github.com/mg98/scriptup/pkg/scriptup/storage"
)

func MigrateDown(cfg *Config, steps int) error {
	s, err := cfg.InitMigrationState()
	if err != nil {
		return err
	}
	defer s.Close()

	migrations, err := s.All(storage.OrderDesc)
	if err != nil {
		return err
	}
	if steps >= 0 {
		migrations = migrations[:steps]
	}

	for _, m := range migrations {
		if err := m.Parse(cfg.Executor); err != nil {
			return err
		}
		if err := m.Down(); err != nil {
			return err
		}
		if err := s.DropLatest(); err != nil {
			return err
		}
	}

	return nil
}
