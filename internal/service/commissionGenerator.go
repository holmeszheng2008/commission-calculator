package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var resultDir string = "result"

type CommissionRpt struct {
	InvoiceNum  string
	AccountNum  string
	BillTo      string
	ProjectCode string
	ProjectName string
	PostingDate string
	Status      string
	Total       string
}

func GenCommisionRpt() error {
	// key being sales
	commisionRptData := make(map[string][]CommissionRpt)
	for invoiceNum, subMap := range RecData {
		for accountNum, total := range subMap {
			if ArData[invoiceNum] == nil {
				continue
			}
			arRecord := ArData[invoiceNum][0]
			billTo := arRecord.BillTo
			postingDate := arRecord.PostingDate
			sales := arRecord.Sales
			status := arRecord.Status
			sTotal := fmt.Sprintf("%f", total)
			var projectCode string
			var projectName string
			if GlData[invoiceNum] == nil {
				projectCode = ""
				projectName = ""
			} else {
				projectCode = GlData[invoiceNum][0].ProjectCode
				projectName = PiData[projectCode]
			}

			commisionRptData[sales] = append(commisionRptData[sales],
				CommissionRpt{invoiceNum, accountNum, billTo, projectCode, projectName, postingDate, status, sTotal})
		}
	}

	err := os.RemoveAll(resultDir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("OS error: Path %s can't be deleted, files inside are being used.", resultDir))
		return err
	}

	time.Sleep(100 * time.Millisecond)
	err = os.Mkdir(resultDir, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("OS error: Path %s can't be created, files inside are being used.", resultDir))
		return err
	}

	for sales, col := range commisionRptData {
		if sales == "" {
			sales = "NO_NAME"
		}

		fileName := fmt.Sprintf("result/%s.csv", strings.ReplaceAll(sales, " ", "_"))
		filePtr, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatalln(fmt.Sprintf("OS error: File %s can't be created.", fileName))
			return err
		}

		csvWriter := csv.NewWriter(filePtr)

		csvWriter.Write([]string{"Invoice #", "Account #", "Bill To", "Project Code", "Project Name", "Posting Date", "Status", "Total"})
		for _, commissionRpt := range col {
			record := []string{commissionRpt.InvoiceNum, commissionRpt.AccountNum, commissionRpt.BillTo, commissionRpt.ProjectCode,
				commissionRpt.ProjectName, commissionRpt.PostingDate, commissionRpt.Status, commissionRpt.Total}
			csvWriter.Write(record)
		}

		csvWriter.Flush()

		filePtr.Close()
	}

	return nil
}
