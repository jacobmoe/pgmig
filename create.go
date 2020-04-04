package pgmig

import (
	"fmt"
	// "github.com/go-pg/pg"
	"github.com/markbates/pkger"
	"log"
	// "os"
	"time"
)

// Create accepts a name and creates up and down migration files.
func (m Migrator) Create(name string) error {
	t := time.Now()
	formattedTime := t.Format("060102150405")

	upFileName := fmt.Sprintf("%s_%s.down.sql", formattedTime, name)
	downFileName := fmt.Sprintf("%s_%s.up.sql", formattedTime, name)

	_, err := pkger.Create(fmt.Sprintf("%s/%s", m.dirPath, upFileName))
	if err != nil {
		return err
	}

	_, err = pkger.Create(fmt.Sprintf("%s/%s", m.dirPath, downFileName))
	if err != nil {
		return err
	}

	log.Printf("created migration files \n- %s\n- %s\n", upFileName, downFileName)

	return nil
}
