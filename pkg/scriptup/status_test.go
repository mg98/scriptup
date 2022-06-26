package scriptup

import (
	"os"
	"testing"
)

func TestGetUnappliedMigrations(t *testing.T) {
	s, err := testCfg.InitMigrationState()
	if err != nil {
		panic(err)
	}
	d := []byte("20220618100116_create-file\n20220618104836_append-1-to-file\n")
	if err := os.WriteFile(testCfg.FileDB, d, 0644); err != nil {
		panic(err)
	}
	defer os.Remove(testCfg.FileDB)

	ms, err := getUnappliedMigrations(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ms) != 2 {
		t.Fatalf("expected 2 unapplied migrations, got %d", len(ms))
	}
	if ms[0].Name != "20220618104839_append-2-to-file" {
		t.Fatalf("expected first to equal 20220618104839_append-2-to-file, got %s", ms[0].Name)
	}
	if ms[1].Name != "20220618104842_append-3-to-file" {
		t.Fatalf("expected first to equal 20220618104842_append-3-to-file, got %s", ms[1].Name)
	}
}
