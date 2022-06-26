package storage

import (
	"fmt"
	"testing"
)

var testPostgresConnDetails = &SQLConnectionDetails{
	Dialect:      "mysql",
	Host:         "localhost",
	Port:         6632,
	User:         "root",
	Password:     "123",
	DatabaseName: "scriptup",
	SSLMode:      "disable",
	TableName:    "migrations",
}

var testMySQLConnDetails = &SQLConnectionDetails{
	Dialect:      "mysql",
	Host:         "localhost",
	Port:         6632,
	User:         "root",
	Password:     "123",
	DatabaseName: "scriptup",
	SSLMode:      "disable",
	TableName:    "migrations",
}

func runTestForEveryDialect(t *testing.T, runTests func(s *SQLStorage)) {
	for _, s := range []*SQLStorage{
		NewSQLStorage(testPostgresConnDetails),
		NewSQLStorage(testMySQLConnDetails),
	} {
		t.Run(fmt.Sprintf("test for dialect %s", s.connDetails.Dialect), func(t *testing.T) {
			runTests(s)
		})
	}
}

func TestSQLStorage_Open(t *testing.T) {
	runTestForEveryDialect(t, func(s *SQLStorage) {
		if err := s.Open(); err != nil {
			t.Fatalf("Open failed: %v", err)
		}
		s.Close()
	})
}

func TestSQLStorage_Close(t *testing.T) {
	runTestForEveryDialect(t, func(s *SQLStorage) {
		if err := s.Open(); err != nil {
			panic(err)
		}
		if err := s.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	})
}

func TestSQLStorage_Append(t *testing.T) {
	runTestForEveryDialect(t, func(s *SQLStorage) {
		if err := s.Open(); err != nil {
			panic(err)
		}

		if err := s.Append("20201224000000_dummy1"); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
		defer s.conn.Exec("DELETE FROM migrations WHERE migration_name = '20201224000000_dummy1'")

		t.Run("db record created", func(t *testing.T) {
			row := s.conn.QueryRow("SELECT * FROM migrations WHERE migration_name = '20201224000000_dummy1'")
			var name string
			if err := row.Scan(&name); err != nil {
				t.Fatalf("record not found: %v", err)
			}
		})
		t.Run("record is unique", func(t *testing.T) {
			if err := s.Append("20201224000000_dummy1"); err == nil {
				t.Fatalf("Append expected to fail")
			}
		})
	})
}

func TestSQLStorage_Pop(t *testing.T) {
	runTestForEveryDialect(t, func(s *SQLStorage) {
		if err := s.Open(); err != nil {
			panic(err)
		}

		if _, err := s.conn.Exec("INSERT INTO migrations VALUES " +
			"('20220619000000-dummy'), " +
			"('20221130000000-dummy'), " +
			"('20210619001000-dummy')"); err != nil {
			panic(err)
		}
		defer s.conn.Exec("DELETE FROM migrations")

		t.Run("pop one migration", func(t *testing.T) {
			if err := s.Pop(); err != nil {
				t.Fatalf("Pop failed: %v", err)
			}

			rows, err := s.conn.Query("SELECT * FROM migrations ORDER BY migration_name ASC")
			if err != nil {
				panic(err)
			}

			var names []string
			for rows.Next() {
				var name string
				if err := rows.Scan(&name); err != nil {
					panic(err)
				}
				names = append(names, name)
			}

			if len(names) != 2 {
				t.Fatalf("expected count of migrations to be 2, got %d", len(names))
			}
			if names[0] != "20210619001000-dummy" {
				t.Fatalf("expected first migration to be %s, got %s", "20210619001000-dummy", names[0])
			}
			if names[1] != "20220619000000-dummy" {
				t.Fatalf("expected first migration to be %s, got %s", "20220619000000-dummy", names[1])
			}
		})

		t.Run("pop remaining migrations", func(t *testing.T) {
			if err := s.Pop(); err != nil {
				t.Fatalf("Pop failed: %v", err)
			}
			if err := s.Pop(); err != nil {
				t.Fatalf("Pop failed: %v", err)
			}
			row := s.conn.QueryRow("SELECT COUNT(*) FROM migrations")
			var count int
			if err := row.Scan(&count); err != nil {
				panic(err)
			}
			if count != 0 {
				t.Fatalf("expected 0 left migrations, found %d", count)
			}
		})

		t.Run("pop one too many", func(t *testing.T) {
			if err := s.Pop(); err == nil {
				t.Fatalf("Pop expected to fail")
			}
		})
	})
}

func TestSQLStorage_All(t *testing.T) {
	runTestForEveryDialect(t, func(s *SQLStorage) {
		if err := s.Open(); err != nil {
			panic(err)
		}
		if _, err := s.conn.Exec("INSERT INTO migrations VALUES " +
			"('20220619000000-dummy'), " +
			"('20221130000000-dummy'), " +
			"('20210619001000-dummy')"); err != nil {
			panic(err)
		}
		defer s.conn.Exec("DELETE FROM migrations")

		entries, err := s.All(OrderAsc)
		if err != nil {
			t.Fatalf("All failed: %v", err)
		}
		if len(entries) != 3 {
			t.Fatalf("expected to find 3 entries, got %d", len(entries))
		}
		if entries[0] != "20210619001000-dummy" {
			t.Fatalf("expected first item to be \"20210619001000-dummy\", got %s", entries[0])
		}
		if entries[1] != "20220619000000-dummy" {
			t.Fatalf("expected first item to be \"20220619000000-dummy\", got %s", entries[1])
		}
		if entries[2] != "20221130000000-dummy" {
			t.Fatalf("expected first item to be \"20221130000000-dummy\", got %s", entries[2])
		}
	})
}
