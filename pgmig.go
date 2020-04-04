package pgmig

import (
	"log"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
)

// Migrator generates and runs migrations
type Migrator struct {
	db      *pg.DB
	dirPath string
}

// New initializes a Migrator
func New(dirPath string) Migrator {
	return Migrator{dirPath: dirPath}
}

// Run accepts a command and runs the migrations.Collection,
// defined through the migration SQL files in the queries dir.
// Commands: up, down, version, set_version [version]
// func Run(db *pg.DB, migrationCmd string) error {
// 	collection := migrations.NewCollection(buildMigrations()...)
// 	collection = collection.DisableSQLAutodiscover(true)

// 	oldVersion, newVersion, err := collection.Run(db, migrationCmd)
// 	if err != nil {
// 		return err
// 	}
// 	if newVersion != oldVersion {
// 		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
// 	} else {
// 		log.Printf("version is %d\n", oldVersion)
// 	}

// 	return nil
// }

func (m Migrator) buildMigrations() []*migrations.Migration {
	queries, err := m.loadMigrationQueries()
	if err != nil {
		panic(err)
	}

	res := []*migrations.Migration{}

	for _, query := range queries {
		res = append(res, &migrations.Migration{
			Version: query.Version,
			UpTx:    true,
			Up:      upMigration(query.Up, query.Version),
			DownTx:  true,
			Down:    downMigration(query.Down, query.Version),
		})
	}

	return res
}

func upMigration(query string, version int64) func(db migrations.DB) error {
	return func(db migrations.DB) error {
		log.Println("running migration", version)
		_, err := db.Exec(query)
		return err
	}

}

func downMigration(query string, version int64) func(db migrations.DB) error {
	return func(db migrations.DB) error {
		log.Println("rolling back migration", version)
		_, err := db.Exec(query)
		return err
	}
}
