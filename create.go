package pgmig

import (
	"fmt"
	"log"
	"os"
	"time"
)

func create(dirPath, name string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("no migrations directory found at %s", dirPath)
	} else if err != nil {
		return fmt.Errorf("problem checking migrations directory at %s: %w", dirPath, err)
	}

	t := time.Now()
	formattedTime := t.Format("060102150405")

	upFileName := fmt.Sprintf("%s_%s.down.sql", formattedTime, name)
	downFileName := fmt.Sprintf("%s_%s.up.sql", formattedTime, name)

	_, err = os.Create(fmt.Sprintf("%s/%s", dirPath, upFileName))
	if err != nil {
		return err
	}

	_, err = os.Create(fmt.Sprintf("%s/%s", dirPath, downFileName))
	if err != nil {
		return err
	}

	log.Printf("created migration files \n- %s\n- %s\n", upFileName, downFileName)

	return nil
}
