package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Ar Invoice No,Posting Date,Bill To Code,Status,Sub Total,Sales Rep,Account
type ArSubStruct struct {
	AccountNo   string
	PostingDate string
	BillTo      string
	Status      string
	Total       float64
	Sales       string
}

// Account,Invoice No,Amount,Proj Code
type GlSubStruct struct {
	AccountNo   string
	Amount      float64
	ProjectCode string
}

// key being invoice number
var ArData map[string][]ArSubStruct
var GlData map[string][]GlSubStruct

// project code -> project name
var PiData map[string]string

// invoice #,account #,total
var RecData map[string]map[string]float64

func findIndex(data []string, target string) int {
	for i, item := range data {
		if strings.TrimSpace(item) == target {
			return i
		}
	}
	return -1
}
func LoadCsvs(arUrlPtr *string, glUrlPtr *string, piUrlPtr *string, recUrlPtr *string) error {
	err := loadAr(*arUrlPtr)
	if err != nil {
		log.Fatalln("AR csv file bad format")
		return err
	}
	// fmt.Println(ArData)

	err = loadGl(*glUrlPtr)
	if err != nil {
		log.Fatalln("GL csv file bad format")
		return err
	}
	// fmt.Println(GlData)

	if piUrlPtr != nil {
		err = loadPi(*piUrlPtr)
		if err != nil {
			log.Fatalln("PI csv file bad format")
			return err
		}
	}
	// fmt.Println(PiData)
	if recUrlPtr != nil {
		err = loadRec(*recUrlPtr)
		if err != nil {
			log.Fatalln("Rec csv file bad format")
			return err
		}
	}
	// fmt.Println(RecData)
	return nil
}

func loadAr(arUrl string) error {
	ArData = make(map[string][]ArSubStruct)
	csvfile, err := os.Open(arUrl)
	if err != nil {
		log.Fatalln("Couldn't open AR csv file", err)
		return err
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.TrimLeadingSpace = true
	isHeader := true
	var invoiceNumIndex int
	var postingDateIndex int
	var billToIndex int
	var statusIndex int
	var subTotalIndex int
	var salesIndex int
	var accountNumIndex int
	for {
		record, err := r.Read()
		// fmt.Println(record)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if isHeader {
			// Ar Invoice No,Posting Date,Bill To Code,Status,Sub Total,Sales Rep,Account
			invoiceNumIndex = findIndex(record, "Ar Invoice No")
			postingDateIndex = findIndex(record, "Posting Date")
			billToIndex = findIndex(record, "Bill To Code")
			statusIndex = findIndex(record, "Status")
			subTotalIndex = findIndex(record, "Sub Total")
			salesIndex = findIndex(record, "Sales Rep")
			accountNumIndex = findIndex(record, "Account")
			if invoiceNumIndex == -1 || postingDateIndex == -1 || billToIndex == -1 || statusIndex == -1 ||
				subTotalIndex == -1 || salesIndex == -1 || accountNumIndex == -1 {
				log.Fatalln("AR csv header bad format")
				return errors.New("AR csv header bad format")
			}
			isHeader = false
			continue
		}
		status := strings.ToUpper(strings.TrimSpace(record[statusIndex]))
		if status != "R" && status != "C" {
			continue
		}

		invoiceNum := record[invoiceNumIndex]
		postingDate := record[postingDateIndex]
		billTo := record[billToIndex]
		sales := record[salesIndex]
		accountNum := record[accountNumIndex]
		total, err := strconv.ParseFloat(record[subTotalIndex], 64)
		if err != nil {
			log.Fatalln(fmt.Sprintf("invoice number %s, account number %s in AR csv has an unparsable amount %s", invoiceNum, accountNum, record[4]))
		}

		subData := ArSubStruct{accountNum, postingDate, billTo, status, total, sales}
		ArData[invoiceNum] = append(ArData[invoiceNum], subData)
	}

	return nil
}

func loadGl(glUrl string) error {
	GlData = make(map[string][]GlSubStruct)
	csvfile, err := os.Open(glUrl)
	if err != nil {
		log.Fatalln("Couldn't open GL csv file", err)
		return err
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.TrimLeadingSpace = true
	isHeader := true
	var invoiceNumIndex int
	var accountNumIndex int
	var projectCodeIndex int
	var amountIndex int

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if isHeader {
			// Account,Invoice No,Amount,Proj Code
			accountNumIndex = findIndex(record, "Account")
			invoiceNumIndex = findIndex(record, "Invoice No")
			amountIndex = findIndex(record, "Amount")
			projectCodeIndex = findIndex(record, "Proj Code")

			if invoiceNumIndex == -1 || accountNumIndex == -1 || amountIndex == -1 || projectCodeIndex == -1 {
				log.Fatalln("GL csv header bad format")
				return errors.New("GL csv header bad format")
			}

			isHeader = false
			continue
		}
		accountNum := record[accountNumIndex]
		invoiceNum := record[invoiceNumIndex]
		projectCode := record[projectCodeIndex]

		f, err := strconv.ParseFloat(record[amountIndex], 64)
		if err != nil {
			log.Fatalln(fmt.Sprintf("invoice number %s, account number %s in GL csv has an unparsable amount %s", invoiceNum, accountNum, record[2]))
		}
		subData := GlSubStruct{accountNum, f, projectCode}

		GlData[invoiceNum] = append(GlData[invoiceNum], subData)
	}

	return nil
}

func loadPi(piUrl string) error {
	PiData = make(map[string]string)
	csvfile, err := os.Open(piUrl)
	if err != nil {
		log.Fatalln("Couldn't open PI csv file", err)
		return err
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.TrimLeadingSpace = true
	isHeader := true
	var projectCodeIndex int
	var projectNameIndex int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if isHeader {
			//Project Code,Project Name
			projectCodeIndex = findIndex(record, "Project Code")
			projectNameIndex = findIndex(record, "Project Name")
			if projectCodeIndex == -1 || projectNameIndex == -1 {
				log.Fatalln("PI csv header bad format")
				return errors.New("PI csv header bad format")
			}
			isHeader = false
			continue
		}

		PiData[record[projectCodeIndex]] = record[projectNameIndex]
	}

	return nil
}

func loadRec(recUrl string) error {
	RecData = make(map[string]map[string]float64)
	csvfile, err := os.Open(recUrl)
	if err != nil {
		log.Fatalln("Couldn't open REC csv file", err)
		return err
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	r.TrimLeadingSpace = true
	isHeader := true
	var invoiceNumIndex int
	var accountNumIndex int
	var totalIndex int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if isHeader {
			// invoice #,account #,total
			invoiceNumIndex = findIndex(record, "invoice #")
			accountNumIndex = findIndex(record, "account #")
			totalIndex = findIndex(record, "total")
			if invoiceNumIndex == -1 || accountNumIndex == -1 || totalIndex == -1 {
				log.Fatalln("REC csv header bad format")
				return errors.New("REC csv header bad format")
			}
			isHeader = false
			continue
		}
		invoiceNum := record[invoiceNumIndex]
		accountNum := record[accountNumIndex]
		total, err := strconv.ParseFloat(record[totalIndex], 64)
		if err != nil {
			log.Fatalln(fmt.Sprintf("invoice number %s, account number %s in REC csv has an unparsable amount %s", invoiceNum, accountNum, record[totalIndex]))
		}
		subMap := RecData[invoiceNum]
		if subMap == nil {
			subMap = make(map[string]float64)
			RecData[invoiceNum] = subMap
		}
		subMap[accountNum] = total
	}

	return nil
}
