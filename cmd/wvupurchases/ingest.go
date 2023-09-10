package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/urfave/cli/v3"
	"github.com/xuri/excelize/v2"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ingest(c *cli.Context) error {
	inFile := c.Args().First()
	if inFile == "" {
		log.Fatal("no input file specified")
	}

	if !(c.String("type") == ProcurementRecordType || c.String("type") == PCardRecordType) {
		log.Fatal("invalid record type specified")
	}

	records, err := xlsxToRecords(inFile, c.String("type"))
	if err != nil {
		log.Fatal("error converting xlsx to records: ", err)
	}

	switch c.String("type") {
	case "procurement":
		var pcpsRecords []ProcurementRecord
		for _, r := range records {
			pcpsRecords = append(pcpsRecords, r.ProcurementRecord)
		}

		err = recordsToDB(pcpsRecords, c.String("db"), c.Bool("overwrite"))
	case "pcard":
		var pcardRecords []PCardRecord
		for _, r := range records {
			pcardRecords = append(pcardRecords, r.PCardRecord)
		}

		err = recordsToDB(pcardRecords, c.String("db"), c.Bool("overwrite"))
	}

	return err
}

func recordsToDB[T ProcurementRecord | PCardRecord](records []T, dbFile string, overwrite bool) error {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	reflect.ValueOf(&records).Elem()

	if err := db.AutoMigrate(new(T)); err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	n := 0
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, record := range records {
			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("error saving record to database: %w", err)
			}
			n++
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error saving records to database: %w", err)
	}

	fmt.Printf("saved %d records to database\n", n)

	return nil
}

func xlsxToRecords(filename string, recordType string) ([]Record, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %w", err)
	}

	defer f.Close()

	sh, err := selectSheet(f)
	if err != nil {
		return nil, fmt.Errorf("error selecting sheet: %w", err)
	}

	rows, err := f.Rows(sh)
	if err != nil {
		return nil, fmt.Errorf("error getting rows from sheet: %w", err)
	}

	var (
		i       = 0
		dec     *RecordDecoder
		records []Record
	)

	for rows.Next() {
		row, _ := rows.Columns()

		if i == 0 {
			dec = NewRecordDecoder(row, recordType)
		} else {
			record, err := dec.UnmarshalRow(row)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling row: %w", err)
			}
			records = append(records, record)
		}

		i += 1
	}

	return records, nil
}
