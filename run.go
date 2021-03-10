package pgmig

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"

	"github.com/jacobmoe/pgmig/errors"
)

const (
	migrationNamePattern = "^([0-9]+)_(.+)\\.(up|down)\\.sql$"
	nameMatchSize        = 4
	indexName            = 0
	indexNum             = 1
	indexDirection       = 3
)

// migrationQuery is an up and down migration with a version num
type migrationQuery struct {
	Up      string
	Down    string
	Version int64
}

func run(db *pg.DB, migFS fs.FS, migrationCmd string) error {
	pgMigrations := buildMigrations(migFS)
	collection := migrations.NewCollection(pgMigrations...)
	collection = collection.DisableSQLAutodiscover(true)

	oldVersion, newVersion, err := collection.Run(db, migrationCmd)
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	return nil
}

func buildMigrations(migFS fs.FS) []*migrations.Migration {
	queries, err := load(migFS)
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

// load loads migration queries from the migrations file system that match migrationNamePattern.
// Results are ordered by migrationQuery version
func load(migFS fs.FS) ([]migrationQuery, error) {
	queryNames := make([]string, 0)
	queryPaths := make(map[string]string)

	res := make([]migrationQuery, 0)

	queryNameReg := regexp.MustCompile(migrationNamePattern)

	err := fs.WalkDir(migFS, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if queryNameReg.Match([]byte(info.Name())) {
			queryNames = append(queryNames, info.Name())
			queryPaths[info.Name()] = path
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	sort.Strings(queryNames)

	for {
		if len(queryNames) == 0 {
			break
		}

		if len(queryNames) < 2 {
			return res, errors.NewSQLLoadErr(
				"migrations must be in up and down pairs",
			)
		}

		currentNames := queryNames[0:2]

		matchGroup1 := queryNameReg.FindStringSubmatch(currentNames[0])
		matchGroup2 := queryNameReg.FindStringSubmatch(currentNames[1])

		if len(matchGroup1) != nameMatchSize || len(matchGroup2) != nameMatchSize {
			message := fmt.Sprintf(
				"migration file names must conform to pattern %s",
				migrationNamePattern,
			)

			return res, errors.NewSQLLoadErr(message)
		}

		var upMigrationName, upMigrationNum, downMigrationName, downMigrationNum string

		if matchGroup1[indexDirection] == "up" && matchGroup2[indexDirection] == "down" {
			upMigrationName = matchGroup1[indexName]
			upMigrationNum = matchGroup1[indexNum]
			downMigrationName = matchGroup2[indexName]
			downMigrationNum = matchGroup2[indexNum]
		} else if matchGroup1[indexDirection] == "down" && matchGroup2[indexDirection] == "up" {
			upMigrationName = matchGroup2[indexName]
			upMigrationNum = matchGroup2[indexNum]
			downMigrationName = matchGroup1[indexName]
			downMigrationNum = matchGroup1[indexNum]
		} else {
			return res, errors.NewSQLLoadErr(
				"migrations must be in up and down pairs",
			)
		}

		if upMigrationNum != downMigrationNum {
			return res, errors.NewSQLLoadErr(
				"migrations must be in matching up and down pairs",
			)
		}

		upMigrationQuery, err := readFile(migFS, queryPaths[upMigrationName])
		if err != nil {
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf("%s %s", upMigrationName, "migration missing"),
			)
		}

		downMigrationQuery, err := readFile(migFS, queryPaths[downMigrationName])
		if err != nil {
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf("%s %s", downMigrationName, "migration missing"),
			)
		}

		version, err := strconv.ParseInt(upMigrationNum, 10, 64)

		if err != nil {
			message := "migration file name must begin with a valid version number"
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf("%s %s", upMigrationName, message),
			)
		}

		res = append(res, migrationQuery{
			Version: version,
			Up:      string(upMigrationQuery),
			Down:    string(downMigrationQuery),
		})

		queryNames = queryNames[2:]
	}

	return res, nil
}

func readFile(migFS fs.FS, path string) (string, error) {
	file, err := migFS.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(file)

	return buf.String(), nil
}
