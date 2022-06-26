package storage

import (
	"os"
	"reflect"
	"testing"
)

const testFilePath = ".test/filedb.txt"

const testMigrationDir = ".test/migrations"

func TestFileStorage_Open(t *testing.T) {
	s := NewFileStorage(testFilePath)

	t.Run("create file if it does not exist", func(t *testing.T) {
		defer os.Remove(testFilePath)
		if err := s.Open(); err != nil {
			t.Fatalf("Open failed: %v", err)
		}
		if _, err := os.Stat(testFilePath); err != nil {
			t.Fatalf("file db file not created: %v", err)
		}
	})

	t.Run("use existing file", func(t *testing.T) {
		if _, err := os.Create(testFilePath); err != nil {
			panic(err)
		}
		defer os.Remove(testFilePath)
		if err := s.Open(); err != nil {
			t.Fatalf("Open failed: %v", err)
		}
	})

}

func TestFileStorage_Close(t *testing.T) {
	if _, err := os.Create(testFilePath); err != nil {
		panic(err)
	}
	defer os.Remove(testFilePath)
	s := NewFileStorage(testFilePath)
	if err := s.Open(); err != nil {
		panic(err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

func TestFileStorage_Append(t *testing.T) {
	if _, err := os.Create(testFilePath); err != nil {
		panic(err)
	}
	defer os.Remove(testFilePath)

	s := NewFileStorage(testFilePath)
	if err := s.Open(); err != nil {
		panic(err)
	}
	{
		if err := s.Append("20200101000000_dummy"); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
		b, err := os.ReadFile(testFilePath)
		if err != nil {
			panic(err)
		}
		if string(b) != "20200101000000_dummy\n" {
			t.Fatalf("file content not as expected, got: %s", string(b))
		}
	}
	{
		if err := s.Append("20200101000001_dummy"); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
		if err := s.Append("20200101000002_dummy"); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
		b, err := os.ReadFile(testFilePath)
		if err != nil {
			panic(err)
		}
		if string(b) != "20200101000000_dummy\n20200101000001_dummy\n20200101000002_dummy\n" {
			t.Fatalf("file content not as expected, got: %s", string(b))
		}
	}
}

func TestFileStorage_Pop(t *testing.T) {
	f, err := os.Create(testFilePath)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testFilePath)
	_, err = f.WriteString("20200101000000_dummy\n20200101000001_dummy\n20200101000002_dummy\n")
	if err != nil {
		panic(err)
	}

	s := NewFileStorage(testFilePath)
	if err := s.Open(); err != nil {
		panic(err)
	}
	{
		if err := s.Pop(); err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		b, err := os.ReadFile(testFilePath)
		if err != nil {
			panic(err)
		}
		if string(b) != "20200101000000_dummy\n20200101000001_dummy\n" {
			t.Fatalf("file content not as expected, got: %s", string(b))
		}
	}
	{
		if err := s.Pop(); err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		b, err := os.ReadFile(testFilePath)
		if err != nil {
			panic(err)
		}
		if string(b) != "20200101000000_dummy\n" {
			t.Fatalf("file content not as expected, got: %s", string(b))
		}
	}
	{
		if err := s.Pop(); err != nil {
			t.Fatalf("Pop failed: %v", err)
		}
		b, err := os.ReadFile(testFilePath)
		if err != nil {
			panic(err)
		}
		if string(b) != "" {
			t.Fatalf("file content not as expected, got: %s", string(b))
		}
	}
	{
		if err := s.Pop(); err.Error() != "no migrations left" {
			t.Fatalf("expected error \"no migrations left\", got: %v", err)
		}
	}

}

func TestFileStorage_All(t *testing.T) {
	f, err := os.Create(testFilePath)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testFilePath)

	{
		s := NewFileStorage(testFilePath)
		if err := s.Open(); err != nil {
			panic(err)
		}
		entries, err := s.All(OrderAsc)
		if err != nil {
			t.Fatalf("All failed: %v", err)
		}

		if len(entries) != 0 {
			t.Fatalf("All expected to be empty, got: %v", entries)
		}
	}

	_, err = f.WriteString("20200101000000_dummy\n20200101000001_dummy\n20200101000002_dummy")
	if err != nil {
		panic(err)
	}

	{
		s := NewFileStorage(testFilePath)
		if err := s.Open(); err != nil {
			panic(err)
		}
		entries, err := s.All(OrderAsc)
		if err != nil {
			t.Fatalf("All failed: %v", err)
		}

		if !reflect.DeepEqual(entries, []string{"20200101000000_dummy", "20200101000001_dummy", "20200101000002_dummy"}) {
			t.Fatalf("All does not return as expected, got: %v", entries)
		}
	}
}

func TestFileStorage_Latest(t *testing.T) {
	t.Run("empty file, i.e., no records", func(t *testing.T) {
		if _, err := os.Create(testFilePath); err != nil {
			panic(err)
		}
		defer os.Remove(testFilePath)

		s := NewFileStorage(testFilePath)
		if err := s.Open(); err != nil {
			panic(err)
		}
		l, err := s.Latest()
		if err != nil {
			t.Fatalf("Latest failed: %v", err)
		}
		if l != nil {
			t.Fatalf("latest expected to be nil, got \"%s\"", *l)
		}
	})

	t.Run("single record", func(t *testing.T) {
		f, err := os.Create(testFilePath)
		if err != nil {
			panic(err)
		}
		defer os.Remove(testFilePath)
		if _, err := f.WriteString("dummy\n"); err != nil {
			panic(err)
		}

		s := NewFileStorage(testFilePath)
		if err := s.Open(); err != nil {
			panic(err)
		}
		l, err := s.Latest()
		if err != nil {
			t.Fatalf("Latest failed: %v", err)
		}
		if *l != "dummy" {
			t.Fatalf("latest expected to be \"dummy\", got \"%s\"", *l)
		}
	})

	t.Run("many records", func(t *testing.T) {
		f, err := os.Create(testFilePath)
		if err != nil {
			panic(err)
		}
		defer os.Remove(testFilePath)
		if _, err := f.WriteString("dummy\ntammy\nsammy\n"); err != nil {
			panic(err)
		}

		s := NewFileStorage(testFilePath)
		if err := s.Open(); err != nil {
			panic(err)
		}
		l, err := s.Latest()
		if err != nil {
			t.Fatalf("Latest failed: %v", err)
		}
		if *l != "sammy" {
			t.Fatalf("latest expected to be \"sammy\", got \"%s\"", *l)
		}
	})
}
