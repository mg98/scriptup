package storage

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type FileStorage struct {
	Path string
	file *os.File
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

// Open the file.
func (s *FileStorage) Open() error {
	var err error
	s.file, err = os.OpenFile(s.Path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	return err
}

// Close the file.
func (s *FileStorage) Close() error {
	return s.file.Close()
}

// Append the entry as a new line at the end of the file.
func (s *FileStorage) Append(entry string) error {
	if err := s.Open(); err != nil {
		return err
	}
	defer s.Close()
	_, err := s.file.WriteString(entry + "\n")
	return err
}

// Pop removes the last line from the file.
func (s *FileStorage) Pop() error {
	if err := s.Open(); err != nil {
		return err
	}
	defer s.Close()
	var lines []string
	scanner := bufio.NewScanner(s.file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(lines) == 0 {
		return errors.New("no migrations left")
	}
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	for _, line := range lines[:len(lines)-1] {
		if _, err := s.file.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// All returns all entries from the file.
func (s *FileStorage) All(o Order) ([]string, error) {
	if err := s.Open(); err != nil {
		return nil, err
	}
	defer s.Close()
	var entries []string
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		entries = append(entries, strings.Trim(scanner.Text(), " \n"))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if o == OrderAsc {
		return entries, nil
	}

	// reverse for descending order
	var reverseEntries []string
	for i := len(entries) - 1; i >= 0; i-- {
		reverseEntries = append(reverseEntries, entries[i])
	}
	return reverseEntries, nil
}

// Latest returns the most recently added entry.
func (s *FileStorage) Latest() (*string, error) {
	if err := s.Open(); err != nil {
		return nil, err
	}
	defer s.Close()
	scanner := bufio.NewScanner(s.file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if lastLine == "" {
		return nil, nil
	}
	return &lastLine, nil
}
