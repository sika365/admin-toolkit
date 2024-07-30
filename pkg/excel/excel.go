package excel

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/sirupsen/logrus"
	"gitlab.sikapp.ir/sikatech/eshop/eshop-sdk-go-v1/models"

	"github.com/sika365/admin-tools/context"
	"github.com/sika365/admin-tools/pkg/file"
	"github.com/sika365/admin-tools/utils"
)

const (
	ExcelRegex = `^.*\.xlsx?$`
	CSVRegex   = `^.*\.csv?$`
)

type ScanRequest struct {
	models.NodeRequest
	file.ScanRequest
	Offset int `json:"offset,omitempty"`
}

type ConvertRequest struct {
	ScanRequest
	OutputPath   string `json:"output_path,omitempty"`
	ForceReplace bool   `json:"force_replace,omitempty"`
}

func ConvertExcelsToCSVs(ctx *context.Context, inputDir, outputDir string, forceReplace bool, excelFiles file.MapFiles) error {
	if inputDir == "" {
		inputDir = "."
	}
	if outputDir == "" {
		outputDir = inputDir
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.FileMode(0775)); err != nil && !os.IsExist(err) {
			return err
		}
	}

	for _, f := range excelFiles {
		xlsxPath := f.Path
		csvPath := path.Join(outputDir, fmt.Sprintf("%s.csv", f.Name))

		if _, err := os.Stat(csvPath); err == nil && forceReplace {
			if err := os.Remove(csvPath); err != nil {
				return err
			}
		}

		if err := utils.ConvertXLSXToCSV(xlsxPath, csvPath); err != nil {
			return err
		}
	}

	return nil
}

func LoadExcels(ctx *context.Context, inputDir string, maxDepth int) (file.MapFiles, error) {
	if inputDir == "" {
		inputDir = "."
	}
	tempDir := inputDir

	if reExcel, err := regexp.Compile(ExcelRegex); err != nil {
		return nil, err
	} else if xlsxFiles, _ := file.WalkDir(inputDir, maxDepth, nil, reExcel); false {
		logrus.Info("!!! no excel files found !!!")
	} else if err := ConvertExcelsToCSVs(ctx, inputDir, tempDir, true, xlsxFiles); err != nil {
		logrus.Info("xxx convert File to csv failed xxx")
	}

	if reCSV, err := regexp.Compile(CSVRegex); err != nil {
		return nil, err
	} else if csvFiles, _ := file.WalkDir(tempDir, maxDepth, nil, reCSV); len(csvFiles) == 0 {
		logrus.Info("!!! no excel files found !!!")
		return nil, nil
	} else {
		return csvFiles, nil
	}
}

func FromFiles(files file.MapFiles, offset int, fn func(header map[string]int, rec []string)) error {
	for _, f := range files {
		reader := csv.NewReader(f.Open().Reader())
		// Read header
		header := make(map[string]int)
		if h, err := reader.Read(); err != nil {
			return fmt.Errorf("failed to read header row %d: %w", 1, err)
		} else {
			for i, t := range h {
				header[t] = i
			}
		}
		// Skip the specified number of rows (offset)
		for i := 1; i < offset; i++ {
			if _, err := reader.Read(); err != nil {
				return fmt.Errorf("failed to skip row %d: %w", i+1, err)
			}
		}
		// Read the remaining records from the CSV file
		i := offset + 1 // 1 row header
		for {
			if r, err := reader.Read(); err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("failed to read row %d: %w", i+1, err)
			} else if errors.Is(err, io.EOF) {
				break
			} else {
				fn(header, r)
			}
		}

		f.Close()
	}
	return nil
}
