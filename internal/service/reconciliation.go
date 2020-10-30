package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

var recFilePath string = "reconciliation.csv"
var tolerance float64 = 0.0001
var ArAccumulatedData map[string]map[string]float64
var hitMemo map[string]map[string]bool = make(map[string]map[string]bool)

func Reconciliate() error {
	glAccumulatedData := make(map[string]map[string]float64)
	mismatches := 0
	err := os.Remove(recFilePath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln(fmt.Sprintf("OS error: File %s being used, can't be deleted", recFilePath))
		return err
	}
	time.Sleep(100 * time.Millisecond)
	reconFilePtr, err := os.OpenFile(recFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(fmt.Sprintf("OS error: File %s can't be created.", recFilePath))
		return err
	}
	defer reconFilePtr.Close()

	ArAccumulatedData = make(map[string]map[string]float64)

	for invoiceNum, col := range ArData {
		ArAccumulatedData[invoiceNum] = make(map[string]float64)
		for _, arSubStruct := range col {
			accountNum := arSubStruct.AccountNo
			amount := arSubStruct.Total
			ArAccumulatedData[invoiceNum][accountNum] = ArAccumulatedData[invoiceNum][accountNum] + amount
		}
	}

	for invoiceNum, col := range GlData {
		glAccumulatedData[invoiceNum] = make(map[string]float64)
		for _, glSubStruct := range col {
			accountNum := glSubStruct.AccountNo
			amount := glSubStruct.Amount
			glAccumulatedData[invoiceNum][accountNum] = glAccumulatedData[invoiceNum][accountNum] + amount
		}
	}

	csvWriter := csv.NewWriter(reconFilePtr)
	csvWriter.Write([]string{"invoice #", "account #", "AR total", "GL total", "match"})
	for invoiceNum, subMap := range ArAccumulatedData {
		hitMemo[invoiceNum] = make(map[string]bool)
		for accountNum, amount := range subMap {
			var glTotal float64 = 0
			glSubMap := glAccumulatedData[invoiceNum]
			if glSubMap != nil {
				glTotal = glSubMap[accountNum]
			}
			match := math.Abs(amount-glTotal) < tolerance
			if !match {
				mismatches++
			}
			record := []string{invoiceNum, accountNum, fmt.Sprintf("%f", amount), fmt.Sprintf("%f", glTotal),
				strconv.FormatBool(match)}
			csvWriter.Write(record)
			hitMemo[invoiceNum][accountNum] = true
		}
	}

	for invoiceNum, subMap := range glAccumulatedData {
		for accountNum, amount := range subMap {
			if hitMemo[invoiceNum][accountNum] {
				continue
			}
			match := math.Abs(amount) < tolerance
			if !match {
				mismatches++
			}
			record := []string{invoiceNum, accountNum, fmt.Sprintf("%f", 0.0), fmt.Sprintf("%f", amount),
				strconv.FormatBool(match)}
			csvWriter.Write(record)
		}
	}

	glAccumulatedData = nil

	csvWriter.Flush()

	if mismatches == 0 {
		return nil
	}

	return fmt.Errorf("Reconciliation failed, %d mismatch(es) occured", mismatches)
}
