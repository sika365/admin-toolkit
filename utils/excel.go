package utils

import (
	"encoding/csv"
	"os"

	"github.com/xuri/excelize/v2"
)

// ConvertXLSXToCSV converts a given XLSX file to CSV
func ConvertXLSXToCSV(xlsxPath string, csvPath string) error {
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create CSV file
	csvFile, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Get all rows in the first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return err
	}

	// Write rows to CSV
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
