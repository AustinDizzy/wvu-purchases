package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ProcurementRecordType = "procurement"
	PCardRecordType       = "pcard"
)

type (
	ProcurementRecord struct {
		PurchasingOrg   string `excel:"Purchasing Org"`
		ApprovedDate    string `excel:"Approved Date"`
		ApprovedTime    string
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

	PCardRecord struct {
		AccountName            string  `excel:"ACCOUNT NAME"`
		PostingDate            string  `excel:"POSTING DATE,to_date"`
		TransDate              string  `excel:"TRANS DATE,to_date"`
		MCC                    string  `excel:"MCC"`
		MCCDescription         string  `excel:"MCC DESCRIPTION"`
		MerchantName           string  `excel:"MERCHANT NAME"`
		TransAmount            float64 `excel:"TRANS AMOUNT,to_num"`
		EBO                    string  `excel:"EBO"`
		ReportsTo              string  `excel:"REPORTS TO INTERMEDIATE"`
		DCC                    string  `excel:"DCC"`
		Identifier             string  `excel:"IDENTIFIER"`
		SalesTaxAmount         float64 `excel:"SALES TAX AMOUNT,to_num"`
		StateProvince          string  `excel:"STATE / PROVINCE"`
		CityName               string  `excel:"CITY NAME"`
		ShipFromZip            string  `excel:"SHIP FROM POSTAL CODE"`
		CountryCode            string  `excel:"COUNTRY CODE"`
		CustomerCode           string  `excel:"CUSTOMER CODE"`
		ExpenseDesc            string  `excel:"EXPENSE DESCRIPTION"`
		ItemDesc               string  `excel:"ITEM DESCRIPTION"`
		GP                     string  `excel:"G/P"`
		Campus                 string  `excel:"CAMPUS"`
		DeptActivity           string  `excel:"DEPARTMENTAL ACTIVITY"`
		Fund                   string  `excel:"FUND"`
		LineItem               string  `excel:"LINE ITEM"`
		Function               string  `excel:"FUNCTION"`
		Project                string  `excel:"PROJECT"`
		Task                   string  `excel:"TASK"`
		Award                  string  `excel:"AWARD"`
		ExpenditureType        string  `excel:"EXPENDITURE TYPE"`
		ExpenditureOrg1        string  `excel:"EXPENDITURE ORG 1"`
		ExpenditureOrg2        string  `excel:"EXPENDITURE ORG 2"`
		AirCityOfDest          string  `excel:"AIR CITY OF DESTINATION"`
		AirCityOfOrigin        string  `excel:"AIR CITY OF ORIGIN"`
		AirStopOverIndicator   string  `excel:"AIR STOP OVER INDICATOR"`
		TicketNumber           string  `excel:"TICKET NUMBER"`
		PassengerName          string  `excel:"PASSENGER NAME"`
		VRCheckoutDate         string  `excel:"VR CHECKOUT DATE"`
		VRRentersName          string  `excel:"VR RENTER'S NAME"`
		LODArrivalDate         string  `excel:"LOD ARRIVAL DATE"`
		LODTotalRoomNights     string  `excel:"LOD TOTAL ROOM NIGHTS"`
		FLTItemDescription     string  `excel:"FLT ITEM DESCRIPTION"`
		FLTMotorFuelAmount     string  `excel:"FLT MOTOR FUEL AMOUNT"`
		FLTMotorFuelSaleAmount string  `excel:"FLT MOTOR FUEL SALE AMOUNT"`
		FLTMotorFuelUnitAmount string  `excel:"FLT MOTOR FUEL UNIT AMOUNT"`
		ItemQuantity           string  `excel:"ITEM QUANTITY"`
		ExtendedAmount         float64 `excel:"EXTENDED AMOUNT,to_num"`
		UnitPrice              float64 `excel:"UNIT PRICE,to_num"`
		StateOrRCCard          string  `excel:"STATE OR RC CARD"`
	}

	Record struct {
		ProcurementRecord
		PCardRecord
	}

	RecordDecoder struct {
		colMap  map[int]string
		decType reflect.Type
	}
)

func (r ProcurementRecord) TableName() string {
	return "procurement_records"
}

func (r PCardRecord) TableName() string {
	return "pcard_records"
}

func NewRecordDecoder(colRow []string, decType string) *RecordDecoder {
	m := map[int]string{}
	for i, col := range colRow {
		m[i] = col
	}

	var t reflect.Type

	switch decType {
	case ProcurementRecordType:
		t = reflect.TypeOf(ProcurementRecord{})
	case PCardRecordType:
		t = reflect.TypeOf(PCardRecord{})
	default:
		panic(fmt.Sprintf("unknown decType: %s", decType))
	}

	return &RecordDecoder{
		colMap:  m,
		decType: t,
	}
}

func (dec *RecordDecoder) UnmarshalRow(r []string) (Record, error) {
	record := Record{}

	for i, v := range r {
		colName, ok := dec.colMap[i]
		if !ok {
			return record, fmt.Errorf("no value for colMap[%d]", i)
		}

		for _, f := range reflect.VisibleFields(dec.decType) {
			tags := strings.Split(f.Tag.Get("excel"), ",")
			if strings.EqualFold(tags[0], colName) {
				if len(tags) == 1 {
					reflect.ValueOf(&record).Elem().FieldByName(f.Name).SetString(v)
				}

				for _, t := range tags[1:] {
					switch t {
					case "to_num":
						v = regexp.MustCompile(`[\s,$]`).ReplaceAllString(v, "")
						if v == "" {
							v = "0"
						}
						numVal, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return record, fmt.Errorf("failed to parse num in to_num (i= %d, field = %s, val = %s): %w", i, colName, v, err)
						}
						reflect.ValueOf(&record).Elem().FieldByName(f.Name).SetFloat(numVal)
					case "to_date":
						v = strings.ReplaceAll(v, "-", "/")
						t, err := time.Parse("01/02/06", v)
						if err != nil {
							return record, fmt.Errorf("failed to parse time in to_date (i= %d, field = %s, val = %s): %w", i, colName, v, err)
						}
						reflect.ValueOf(&record).Elem().FieldByName(f.Name).SetString(t.Format("2006-01-02"))
					}
				}

				if colName == "Approved Date" {
					t, err := time.Parse("1/2/06 15:04", v)
					if err != nil {
						return record, fmt.Errorf("failed to parse time: %w", err)
					}
					record.ApprovedDate = t.Format("2006-01-02")
					record.ApprovedTime = t.Format("15:04")
				}

				break
			}
		}
	}

	return record, nil
}
