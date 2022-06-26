package scriptup

import (
	"fmt"
	"github.com/mg98/scriptup/pkg/scriptup/migration"
	"github.com/mg98/scriptup/pkg/scriptup/migration_state"
)

func Status(cfg *Config) error {
	s, err := cfg.InitMigrationState()
	if err != nil {
		return err
	}
	ms, err := getUnappliedMigrations(s)
	if err != nil {
		return err
	}
	for _, m := range ms {
		fmt.Println(m.Name)
	}

	return nil
}

// getUnappliedMigrations returns all migrations from the migration directory
// which are dated after the latest migration in the storage.
func getUnappliedMigrations(state *migration_state.MigrationState) ([]*migration.Migration, error) {
	latestM, err := state.Latest()
	if err != nil {
		return nil, err
	}
	ms, err := state.Files()
	if err != nil {
		return nil, err
	}
	var res []*migration.Migration
	for _, m := range ms {
		if latestM == nil || m.Date.After(latestM.Date) {
			res = append(res, m)
		}
	}
	return res, nil
}
