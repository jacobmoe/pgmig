package pgmig

import (
	"fmt"
	"log"
	"time"

	"github.com/markbates/pkger"
)

func create(dirPath, name string) error {
	t := time.Now()
	formattedTime := t.Format("060102150405")

	upFileName := fmt.Sprintf("%s_%s.down.sql", formattedTime, name)
	downFileName := fmt.Sprintf("%s_%s.up.sql", formattedTime, name)

	_, err := pkger.Create(fmt.Sprintf("%s/%s", dirPath, upFileName))
	if err != nil {
		return err
	}

	_, err = pkger.Create(fmt.Sprintf("%s/%s", dirPath, downFileName))
	if err != nil {
		return err
	}

	log.Printf("created migration files \n- %s\n- %s\n", upFileName, downFileName)

	return nil
}
