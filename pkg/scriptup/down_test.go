package scriptup

import (
	"errors"
	"os"
	"testing"
)

func TestMigrateDown(t *testing.T) {
	if err := os.WriteFile(createdTestFile, []byte("hello world\n1\n2\n3\n"), 0644); err != nil {
		panic(err)
	}
	defer os.Remove(createdTestFile)
	if err := os.WriteFile(
		testCfg.FileDB,
		[]byte("20220618100116_create-file\n20220618104836_append-1-to-file\n20220618104839_append-2-to-file\n20220618104842_append-3-to-file\n"),
		os.ModePerm,
	); err != nil {
		panic(err)
	}
	defer os.Remove(testCfg.FileDB)

	// down --steps=2
	t.Run("2 steps", func(t *testing.T) {
		if err := MigrateDown(testCfg, 2); err != nil {
			t.Fatalf("command failed: %v", err)
		}
		fileBytes, err := os.ReadFile(createdTestFile)
		if err != nil {
			t.Fatalf("file should still exist after 2 steps")
		}
		if string(fileBytes) != "hello world\n1\n" {
			t.Fatalf("migration ineffective: file does not have expected contents: %s", string(fileBytes))
		}
		storageBytes, err := os.ReadFile(testCfg.FileDB)
		if err != nil {
			panic(err)
		}
		if string(storageBytes) != "20220618100116_create-file\n20220618104836_append-1-to-file\n" {
			t.Fatalf("state file not as expected, got %s", string(storageBytes))
		}
	})

	// down
	t.Run("remaining", func(t *testing.T) {
		if err := MigrateDown(testCfg, -1); err != nil {
			t.Fatalf("command failed: %v", err)
		}
		if _, err := os.Stat(createdTestFile); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("migration ineffective: %v", err)
		}
	})

}
