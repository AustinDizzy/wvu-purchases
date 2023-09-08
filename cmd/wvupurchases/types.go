package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	ProcurementRecord struct {
		PurchasingOrg   string  `excel:"Purchasing Org"`
		Approved        string  `excel:"Approved Date"`
		SupplierContact string  `excel:"Supplier Contact"`
		Number          string  `excel:"Number"`
		Line            string  `excel:"Line"`
		Category        string  `excel:"Category"`
		Description     string  `excel:"Description"`
		Quantity        float64 `excel:"Quantity,to_num"`
		UOM             string  `excel:"UOM"`
		Price           float64 `excel:"Price,to_num"`
		Amount          float64 `excel:"Amount,to_num"`
		Buyer           string  `excel:"Buyer"`
		Cancelled       string  `excel:"Cancelled"`
		ClosureStatus   string  `excel:"Closure Status"`
		MiscAttr        string  `excel:"[ ]"`
		Type            string  `excel:"Type"`
		Supplier        string  `excel:"Supplier"`
		Site            string  `excel:"Site"`
		Currency        string  `excel:"Currency"`
		ItemRev         string  `excel:"Item Rev"`
		ContractNum     string  `excel:"Contract Num"`
		Ver             string  `excel:"Ver"`
		SupplierItem    string  `excel:"Supplier Item"`
		Item            string  `excel:"Item"`
	}

	RecordDecoder struct {
		colMap map[int]string
	}
)

func (r ProcurementRecord) TableName() string {
	return "procurement_records"
}

func NewRecordDecoder(colRow []string) *RecordDecoder {
	m := map[int]string{}
	for i, col := range colRow {
		m[i] = col
	}

	return &RecordDecoder{
		colMap: m,
	}
}

func (dec *RecordDecoder) UnmarshalRow(r []string) (ProcurementRecord, error) {
	record := ProcurementRecord{}

	for i, v := range r {
		colName, ok := dec.colMap[i]
		if !ok {
			return record, fmt.Errorf("no value for colMap[%d]", i)
		}

		for _, f := range reflect.VisibleFields(reflect.TypeOf(ProcurementRecord{})) {
			tags := strings.Split(f.Tag.Get("excel"), ",")
			if tags[0] == colName {
				if len(tags) == 1 {
					reflect.ValueOf(&record).Elem().FieldByName(f.Name).SetString(v)
				}

				for _, t := range tags[1:] {
					if t == "to_num" {
						numVal, err := strconv.ParseFloat(strings.Replace(v, ",", "", -1), 64)
						if err != nil {
							return record, fmt.Errorf("failed to parse num in to_num: %w", err)
						}
						reflect.ValueOf(&record).Elem().FieldByName(f.Name).SetFloat(numVal)
					}
				}
				break
			}
		}
	}

	return record, nil
}
