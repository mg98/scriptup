package migration_state

import (
	"github.com/mg98/scriptup/pkg/scriptup/migration"
	"github.com/mg98/scriptup/pkg/scriptup/storage"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func init() {
	// set working directory to project root
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestMigrationState_Add(t *testing.T) {
	ms := &MigrationState{
		migrationDir: "",
		storage:      &storage.FileStorage{Path: ".test/filedb.txt"},
	}
	defer os.Remove(".test/filedb.txt")

	t.Run("add first migration", func(t *testing.T) {
		{
			err := ms.Add(&migration.Migration{
				Name: "hello",
				Date: time.Now(),
			})
			if err != nil {
				t.Fatalf("Add failed with error: %v", err)
			}

			items, err := ms.storage.All(storage.OrderAsc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(items) != 1 {
				t.Fatalf("storage expected to have 1 item, got %d", len(items))
			}
			if items[0] != "hello" {
				t.Fatalf("item expected to equal hello, got %s", items[0])
			}
		}

		t.Run("add second migration", func(t *testing.T) {
			err := ms.Add(&migration.Migration{
				Name: "world",
				Date: time.Now(),
			})
			if err != nil {
				t.Fatalf("Add failed with error: %v", err)
			}

			items, err := ms.storage.All(storage.OrderAsc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(items) != 2 {
				t.Fatalf("storage expected to have 2 items, got %d", len(items))
			}
			if items[0] != "hello" {
				t.Fatalf("first item expected to equal hello, got %s", items[0])
			}
			if items[1] != "world" {
				t.Fatalf("second item expected to equal world, got %s", items[1])
			}
		})
	})
}

func TestMigrationState_All(t *testing.T) {
	ms := &MigrationState{
		migrationDir: "",
		storage:      &storage.FileStorage{Path: ".test/filedb.txt"},
	}
	defer os.Remove(".test/filedb.txt")
	d := []byte("2020010101000000_hey\n2021010101000000_jude\n2022010101000000_dont\n2023010101000000_make\n2024010101000000_it\n2025010101000000_bad\n")
	if err := os.WriteFile(".test/filedb.txt", d, 0644); err != nil {
		panic(err)
	}

	t.Run("ascending order", func(t *testing.T) {
		items, err := ms.All(storage.OrderAsc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 6 {
			t.Fatalf("expected to get 6 items, got %d", len(items))
		}
		if items[0].Name != "2020010101000000_hey" {
			t.Fatalf("first item expected to equal 2020010101000000_hey, got %s", items[0].Name)
		}
		if items[5].Name != "2025010101000000_bad" {
			t.Fatalf("last item expected to equal 2020010101000000_bad, got %s", items[5].Name)
		}
	})

	t.Run("descending order", func(t *testing.T) {
		items, err := ms.All(storage.OrderDesc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 6 {
			t.Fatalf("expected to get 6 items, got %d", len(items))
		}
		if items[0].Name != "2025010101000000_bad" {
			t.Fatalf("first item expected to equal 2020010101000000_bad, got %s", items[0].Name)
		}
		if items[5].Name != "2020010101000000_hey" {
			t.Fatalf("last item expected to equal 2020010101000000_hey, got %s", items[5].Name)
		}
	})
}

func TestMigrationState_Latest(t *testing.T) {
	ms := &MigrationState{
		migrationDir: "",
		storage:      &storage.FileStorage{Path: ".test/filedb.txt"},
	}
	defer os.Remove(".test/filedb.txt")

	t.Run("get latest of no migrations", func(t *testing.T) {
		m, err := ms.Latest()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m != nil {
			t.Fatalf("latest migration expected to be nil, got %v", m)
		}
	})

	t.Run("get latest of a set of migrations", func(t *testing.T) {
		d := []byte("2020010101000000_hey\n2021010101000000_jude\n2022010101000000_dont\n2023010101000000_make\n2024010101000000_it\n2025010101000000_bad\n")
		if err := os.WriteFile(".test/filedb.txt", d, 0644); err != nil {
			panic(err)
		}

		m, err := ms.Latest()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.Name != "2025010101000000_bad" {
			t.Fatalf("Latest expected to be 2025010101000000_bad, got %s", m.Name)
		}
	})
}
