package scriptup

import (
	"errors"
	"github.com/mg98/scriptup/pkg/scriptup/migration_state"
	"github.com/mg98/scriptup/pkg/scriptup/storage"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	FileDB       string `yaml:"file_db,omitempty"`
	Dialect      string `yaml:"dialect"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"pass,omitempty"`
	SSLMode      string `yaml:"sslmode,omitempty"`
	DatabaseName string `yaml:"db_name"`
	Table        string `yaml:"table"`
	Directory    string `yaml:"dir"`
	Executor     string `yaml:"executor"`
}

// GetConfig loads the config for the given environment from the yaml file.
func GetConfig(env string) *Config {
	yamlFile, err := os.ReadFile("./scriptup.yaml")
	if errors.Is(err, os.ErrNotExist) {
		yamlFile, err = os.ReadFile("./scriptup.yml")
	}
	if err != nil {
		log.Fatal(err)
	}
	cfg := map[string]*Config{}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return cfg[env]
}

// InitMigrationState sets up the storage and the MigrationState according to the config.
func (cfg *Config) InitMigrationState() (*migration_state.MigrationState, error) {
	var s storage.Storage
	if cfg.FileDB != "" {
		s = storage.NewFileStorage(cfg.FileDB)
	} else {
		s = storage.NewSQLStorage(&storage.SQLConnectionDetails{
			Dialect:      cfg.Dialect,
			Host:         cfg.Host,
			Port:         cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			DatabaseName: cfg.DatabaseName,
			TableName:    cfg.Table,
			SSLMode:      cfg.SSLMode,
		})
		if err := s.Open(); err != nil {
			return nil, err
		}
	}

	return migration_state.New(cfg.Directory, s), nil
}

// InitStorage retrieves and opens the storage according to the configuration.
func (cfg *Config) InitStorage() (storage.Storage, error) {
	var s storage.Storage
	if cfg.FileDB != "" {
		s = storage.NewFileStorage(cfg.FileDB)
	} else {
		s = storage.NewSQLStorage(&storage.SQLConnectionDetails{
			Dialect:      cfg.Dialect,
			Host:         cfg.Host,
			Port:         cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			DatabaseName: cfg.DatabaseName,
			TableName:    cfg.Table,
			SSLMode:      cfg.SSLMode,
		})
		if err := s.Open(); err != nil {
			return nil, err
		}
	}
	return s, nil
}
