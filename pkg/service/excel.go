package service

import (
	"chip-distribution-go/pkg/entity"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

const (
	ExcelFileDir = "./data"
	SheetName    = "历史行情"
)

func init() {
	if err := os.RemoveAll(ExcelFileDir); err != nil {
		logrus.Error("remove ExcelFilePath error: %v", err)
	}

	if err := os.Mkdir(ExcelFileDir, os.ModePerm); err != nil {
		logrus.Error("mkdir ExcelFilePath error: %v", err)
	}
}

func readExcelData(filePath string) ([]entity.StockInfo, map[string]struct{}, error) {
	var companyMap = make(map[string]struct{})

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		logrus.Error("open excel file error: %v, filePath: %s", err, filePath)
		return []entity.StockInfo{}, companyMap, err
	}
	defer file.Close()

	rows, err := file.GetRows(SheetName)
	if err != nil {
		logrus.Error("get excel rows data error: %v, filePath: %s", err, filePath)
		return []entity.StockInfo{}, companyMap, err
	}

	if len(rows) < 5 {
		return []entity.StockInfo{}, companyMap, errors.New("excel data is tow little")
	}

	// 去除第一行的空白行
	rows = rows[2:]
	var excelList = make([]entity.StockInfo, 0, len(rows))
	for _, row := range rows {
		if len(row) < 12 {
			continue
		}

		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			continue
		}

		open, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			continue
		}

		close, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			continue
		}

		high, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			continue
		}

		low, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			continue
		}

		avg, err := strconv.ParseFloat(row[7], 64)
		if err != nil {
			continue
		}

		quote, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			continue
		}

		vloume, err := strconv.ParseFloat(row[9], 64)
		if err != nil {
			continue
		}

		money, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			continue
		}

		turnoverRate, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			continue
		}

		excelList = append(excelList, entity.StockInfo{
			Date:          date,
			Code:          row[1],
			Name:          row[2],
			Open:          open,
			Close:         close,
			Highest:       high,
			Lowest:        low,
			Average:       avg,
			QuoteChange:   quote,
			Vloume:        vloume,
			Money:         money,
			TurnoverRate:  turnoverRate,
			Concentration: 0.0,
		})

		companyMap[row[2]] = struct{}{}
	}

	return excelList, companyMap, nil
}
