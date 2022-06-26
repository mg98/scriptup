package scriptup

import (
	"errors"
	"fmt"
	"github.com/mg98/scriptup/pkg/scriptup/migration"
	"os"
	"regexp"
	"time"
)

// templateContent is the initial content in a new migration file.
var templateContent = `### migrate up ###

### migrate down ###
`

// currentTime returns the current time as it would be used inside the migration's identifier.
// This is abstracted into a variable to let it be overridden in testing.
var currentTime = time.Now().UTC

// NewMigrationFile creates a new migration file.
func NewMigrationFile(cfg *Config, name string) error {
	if name == "" {
		return errors.New("no migration name")
	}
	if !regexp.MustCompile(fmt.Sprintf("^%s$", migration.NamePattern)).MatchString(name) {
		return errors.New("name contains illegal characters")
	}
	ts := currentTime().Format(migration.DateLayout)
	return os.WriteFile(
		fmt.Sprintf("%s/%s_%s.sh", cfg.Directory, ts, name),
		[]byte(templateContent),
		0644,
	)
}
