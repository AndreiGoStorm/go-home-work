package sqlstorage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

type Migration struct {
	migrationDir string
}

func NewMigration() *Migration {
	return &Migration{"migrations"}
}

func (m *Migration) migrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrations failed for goose set dialect: %w", err)
	}

	if err := m.searchMigrationDir(); err != nil {
		return fmt.Errorf("migrations failed for wrong migration dir: %w", err)
	}

	if err := goose.Up(db, m.migrationDir); err != nil {
		return fmt.Errorf("migrations failed for goose up: %w", err)
	}
	return nil
}

func (m *Migration) searchMigrationDir() error {
	file, err := filepath.Abs(m.migrationDir)
	if err != nil {
		return err
	}

	if m.dirExists(file) {
		return nil
	}

	file, err = filepath.Abs(".")
	if err != nil {
		return err
	}

	file += "/../" + m.migrationDir
	if m.dirExists(file) {
		m.migrationDir = file
		return nil
	}

	return errors.New("migration dir not found")
}

func (m *Migration) dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
