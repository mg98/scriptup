package scriptup

import (
	"os"
	"testing"
)

func TestMigrateUp(t *testing.T) {
	// up --steps=2
	t.Run("2 steps", func(t *testing.T) {
		if err := MigrateUp(testCfg, 2); err != nil {
			t.Fatalf("command failed: %v", err)
		}
		defer os.Remove(testCfg.FileDB)
		defer os.Remove(createdTestFile)
		b, err := os.ReadFile(createdTestFile)
		if err != nil {
			t.Fatalf("migration ineffective: %v", err)
		}
		if string(b) != "hello world\n1\n" {
			t.Fatalf("migration ineffective: file does not have expected contents: %s", string(b))
		}
	})

	// up
	t.Run("all remaining", func(t *testing.T) {
		if err := MigrateUp(testCfg, -1); err != nil {
			t.Fatalf("command failed: %v", err)
		}
		defer os.Remove(testCfg.FileDB)
		defer os.Remove(createdTestFile)
		b, err := os.ReadFile(createdTestFile)
		if err != nil {
			t.Fatalf("migration ineffective: %v", err)
		}
		if string(b) != "hello world\n1\n2\n3\n" {
			t.Fatalf("migration ineffective: file does not have expected contents: %s", string(b))
		}
	})
}
