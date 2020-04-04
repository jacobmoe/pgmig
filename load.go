package pgmig

import (
	"bytes"
	"fmt"
	"github.com/markbates/pkger"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/jacobmoe/pgmig/errors"
)

const (
	migrationNamePattern = "^([0-9]+)_(.+)\\.(up|down)\\.sql$"
	nameMatchSize        = 4
	indexName            = 0
	indexNum             = 1
	indexDirection       = 3
)

// MigrationQuery is an up and down migration with a version num
type MigrationQuery struct {
	Up      string
	Down    string
	Version int64
}

// loadMigrationQueries loads migration queries from the migrations directory
// that match migrationNamePattern.
// Results are ordered by MigrationQuery version
func (m Migrator) loadMigrationQueries() ([]MigrationQuery, error) {
	queryNames := make([]string, 0)
	queryPaths := make(map[string]string)

	res := make([]MigrationQuery, 0)

	queryNameReg := regexp.MustCompile(migrationNamePattern)

	err := pkger.Walk(m.dirPath, func(path string, info os.FileInfo, err error) error {
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

		upMigrationQuery, err := readFile(queryPaths[upMigrationName])
		if err != nil {
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf(
					"%s %s",
					upMigrationName,
					"migration missing from compiled assets",
				),
			)
		}

		downMigrationQuery, err := readFile(queryPaths[downMigrationName])
		if err != nil {
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf(
					"%s %s",
					downMigrationName,
					"migration missing from compiled assets",
				),
			)
		}

		version, err := strconv.ParseInt(upMigrationNum, 10, 64)

		if err != nil {
			message := "migration file name must begin with a valid version number"
			return res, errors.NewSQLLoadErr(
				fmt.Sprintf("%s %s", upMigrationName, message),
			)
		}

		res = append(res, MigrationQuery{
			Version: version,
			Up:      string(upMigrationQuery),
			Down:    string(downMigrationQuery),
		})

		queryNames = queryNames[2:]
	}

	return res, nil
}

func readFile(path string) (string, error) {
	file, err := pkger.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(file)

	return buf.String(), nil
}
