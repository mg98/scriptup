package storage

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"regexp"
)

type SQLConnectionDetails struct {
	Dialect      string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	TableName    string
	SSLMode      string
}

type SQLStorage struct {
	connDetails *SQLConnectionDetails
	conn        *sql.DB
}

func NewSQLStorage(scd *SQLConnectionDetails) *SQLStorage {
	return &SQLStorage{connDetails: scd}
}

func (s *SQLStorage) Open() error {
	var connString string
	switch s.connDetails.Dialect {
	case "postgres":
		connString = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			s.connDetails.Host,
			s.connDetails.Port,
			s.connDetails.User,
			s.connDetails.Password,
			s.connDetails.DatabaseName,
			s.connDetails.SSLMode,
		)
		break
	case "mysql":
		connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			s.connDetails.User,
			s.connDetails.Password,
			s.connDetails.Host,
			s.connDetails.Port,
			s.connDetails.DatabaseName,
		)
	default:
		return errors.New("sql dialect not supported")
	}

	var err error
	s.conn, err = sql.Open(s.connDetails.Dialect, connString)
	if err != nil {
		return err
	}
	err = s.conn.Ping()
	if err != nil {
		return err
	}
	if err := s.setUp(); err != nil {
		return err
	}
	return nil
}

// Close connection to the database.
func (s *SQLStorage) Close() error {
	return s.conn.Close()
}

// setUp adds the required schema to the database if it does not already exist.
func (s *SQLStorage) setUp() error {
	// loose validation just to protect against sql injection
	if !regexp.MustCompile("^[a-zA-Z0-9-_.]*$").MatchString(s.connDetails.TableName) {
		return errors.New("table name contains invalid characters")
	}
	var query string
	switch s.connDetails.Dialect {
	case "postgres":
		query = fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s (migration_name VARCHAR PRIMARY KEY)",
			s.connDetails.TableName,
		)
		break
	case "mysql":
		query = fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s (migration_name VARCHAR(255), PRIMARY KEY (migration_name))",
			s.connDetails.TableName,
		)
	}
	_, err := s.conn.Exec(query)
	return err
}

// Append inserts a new record for that entry into the database.
func (s *SQLStorage) Append(entry string) error {
	if _, err := s.conn.Exec(fmt.Sprintf("INSERT INTO %s VALUES (?)", s.connDetails.TableName), entry); err != nil {
		return err
	}
	return nil
}

// Pop deletes the most recent record (according to the migration date).
func (s *SQLStorage) Pop() error {
	var query string
	switch s.connDetails.Dialect {
	case "postgres":
		query = fmt.Sprintf(
			"DELETE FROM %s WHERE ctid IN (SELECT ctid FROM %s ORDER BY migration_name DESC LIMIT 1)",
			s.connDetails.TableName, s.connDetails.TableName,
		)
		break
	case "mysql":
		query = fmt.Sprintf("DELETE FROM %s ORDER BY migration_name DESC LIMIT 1", s.connDetails.TableName)
	}
	res, err := s.conn.Exec(query)
	if err != nil {
		return err
	}
	i, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("no migrations left")
	}
	return nil
}

// All returns all entries from the database.
func (s *SQLStorage) All(o Order) ([]string, error) {
	rows, err := s.conn.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY migration_name %s", s.connDetails.TableName, o))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		res = append(res, value)
	}
	return res, nil
}

// Latest returns the most recent record (according to the migration date).
func (s *SQLStorage) Latest() (*string, error) {
	row := s.conn.QueryRow(fmt.Sprintf("SELECT * FROM %s ORDER BY migration_name DESC LIMIT 1", s.connDetails.TableName))
	if err := row.Err(); err != nil {
		return nil, err
	}
	var name string
	if err := row.Scan(&name); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if name == "" {
		return nil, nil
	}
	return &name, nil
}
