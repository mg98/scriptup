package scriptup

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewMigration(t *testing.T) {
	const name = "awesome-migration"
	currentTime = func() time.Time {
		return time.Date(2022, time.May, 4, 4, 20, 1, 0, time.UTC)
	}
	if err := NewMigrationFile(testCfg, name); err != nil {
		t.Fatalf("command failed: %v", err)
	}
	migFile := fmt.Sprintf("%s/20220504042001_%s.sh", testCfg.Directory, name)
	defer os.Remove(migFile)
	if _, err := os.Stat(migFile); err != nil && errors.Is(err, os.ErrNotExist) {
		t.Fatalf("migration file not created")
	} else if err != nil {
		panic(err)
	}
}
