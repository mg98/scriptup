package migration

import (
	"errors"
	"strings"
)

var (
	ErrWrongOrder       = errors.New("migrate up section must appear before migrate down section")
	ErrMissingUpSection = errors.New("missing up section in migration file")
)

// migrateUpHeader identifies the start of the up-script in the migration.
const migrateUpHeader = "### migrate up ###"

// migrateDownHeader identifies the start of the down-script in the migration.
const migrateDownHeader = "### migrate down ###"

// parse returns the up and down commands extracted from a migration file's body
func parse(body string) (up string, down string, err error) {
	upPos := strings.Index(body, migrateUpHeader)
	downPos := strings.Index(body, migrateDownHeader)
	if upPos == -1 {
		return "", "", nil
	}
	if downPos != -1 && downPos < upPos {
		return "", "", ErrWrongOrder
	}
	if downPos != -1 {
		up = body[upPos+len(migrateUpHeader) : downPos]
		down = body[downPos+len(migrateDownHeader):]
	} else {
		up = body[upPos+len(migrateUpHeader):]
	}
	return
}
