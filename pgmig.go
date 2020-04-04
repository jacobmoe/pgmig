package pgmig

import (
	"github.com/go-pg/pg/v9"
)

// Migrator generates and runs migrations
type Migrator struct {
	db      *pg.DB
	dirPath string
}

// New initializes a Migrator
func New(db *pg.DB, dirPath string) Migrator {
	return Migrator{db: db, dirPath: dirPath}
}

// Create accepts a name and creates up and down migration files.
func (m Migrator) Create(name string) error {
	return create(m.dirPath, name)
}

// Run accepts a command and runs the migrations Collection,
// defined through the migration SQL files in the queries dir.
// Commands: init, up, down, version, set_version [version]
func (m Migrator) Run(migrationCmd string) error {
	return run(m.db, m.dirPath, migrationCmd)
}

// Init creates version info table in the database
func (m Migrator) Init() error {
	return m.Run("init")
}

// Migrate runs all new migrations
func (m Migrator) Migrate() error {
	return m.Run("up")
}

// Rollback undoes the latest migration
func (m Migrator) Rollback() error {
	return m.Run("down")
}
