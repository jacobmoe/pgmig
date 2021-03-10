package pgmig

import (
	"github.com/go-pg/pg/v10"
	"io/fs"
)

// Migrator generates and runs migrations
type Migrator struct {
	db    *pg.DB
	migFS fs.FS
}

// New initializes a Migrator with a go-pg connection and
// migrations FS
func New(db *pg.DB, migFS fs.FS) Migrator {
	return Migrator{db: db, migFS: migFS}
}

// Create accepts a name and creates up and down migration files.
func (m Migrator) Create(migDirPath, name string) error {
	return create(migDirPath, name)
}

// Run accepts a command and runs the migrations Collection,
// defined through the migration SQL files in the queries dir.
// Commands: init, up, down, version, set_version [version]
func (m Migrator) Run(migrationCmd string) error {
	return run(m.db, m.migFS, migrationCmd)
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
