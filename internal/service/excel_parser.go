package service

import (
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"

	"project/internal/models"
)

type ExcelParser struct {
	service *HierarchyService
}

func NewExcelParser(service *HierarchyService) *ExcelParser {
	return &ExcelParser{service: service}
}

func (p *ExcelParser) ParseExcelFile(file multipart.File) ([]models.UploadRow, error) {
	if file == nil {
		return nil, NewValidationError("no file provided")
	}
	
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, NewValidationError("invalid excel file: " + err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			// Handle any panics from excelize
		}
	}()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, NewValidationError("no sheet found in excel file")
	}

	rawRows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, NewValidationError("unable to read excel rows: " + err.Error())
	}

	if len(rawRows) < 2 {
		return nil, NewValidationError("excel file has no valid data (needs at least header + 1 data row)")
	}

	return p.parseRows(rawRows[0], rawRows[1:]), nil
}

func (p *ExcelParser) parseRows(headers []string, dataRows [][]string) []models.UploadRow {
	var result []models.UploadRow

	for _, row := range dataRows {
		parsedRow := p.parseSingleRow(headers, row)
		if len(parsedRow.Levels) > 0 {
			result = append(result, parsedRow)
		}
	}

	return result
}

func (p *ExcelParser) parseSingleRow(headers, row []string) models.UploadRow {
	metadataHeaders := map[string]bool{
		"sl no":     true,
		"capacity":  true,
		"mf":        true,
		"latitude":  true,
		"longitude": true,
	}

	var levels []models.LevelData
	extraData := make(map[string]string)

	for i, header := range headers {
		header = strings.TrimSpace(header)
		value := getCellValue(row, i)
		normalized := strings.ToLower(header)

		if metadataHeaders[normalized] && value != "" {
			extraData[header] = value
			continue
		}

		if value != "" {
			levels = append(levels, models.LevelData{
				Header: header,
				Value:  value,
			})
		}
	}

	return models.UploadRow{
		Levels:    levels,
		ExtraData: extraData,
	}
}
