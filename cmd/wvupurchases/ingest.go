package main

import (
	"fmt"
	"log"

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

	if !(c.String("type") == "procurement" || c.String("type") == "pcard") {
		log.Fatal("invalid record type specified")
	}

	var err error
	switch c.String("type") {
	case "procurement":
		var records []ProcurementRecord
		records, err = xlsxToRecords(inFile)
		if err != nil {
			log.Fatal("error converting xlsx to records: ", err)
		}
		err = recordsToDB(records, c.String("db"), c.Bool("overwrite"))
	case "pcard":
		log.Println("TODO: pcard")
	}

	return err
}

func recordsToDB(records []ProcurementRecord, dbFile string, overwrite bool) error {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := db.AutoMigrate(&ProcurementRecord{}); err != nil {
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

func xlsxToRecords(filename string) ([]ProcurementRecord, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening input file: %w", err)
	}

	defer f.Close()

	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("error getting rows from sheet: %w", err)
	}

	var (
		i       = 0
		dec     *RecordDecoder
		records []ProcurementRecord
	)

	for rows.Next() {
		row, _ := rows.Columns()

		if i == 0 {
			dec = NewRecordDecoder(row)
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
