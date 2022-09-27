package database

import (
	"chip-distribution-go/pkg/entity"
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DataBaseDir         = "./DB"
	StockDataBasePath   = "./DB/stock.db"
	CompanyDataBasePath = "./DB/company.db"
)

type DataBaseEngine struct {
	stockDB   *gorm.DB
	stockRW   sync.RWMutex
	companyDB *gorm.DB
	companyRW sync.RWMutex
}

func init() {
	_, err := os.Stat(DataBaseDir)
	if os.IsNotExist(err) {
		os.Mkdir(DataBaseDir, os.ModePerm)
	}
}

func NewDataBaseEngine() (*DataBaseEngine, error) {
	stock, err := gorm.Open(sqlite.Open(StockDataBasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	stock.AutoMigrate(&entity.StockInfo{})

	company, err := gorm.Open(sqlite.Open(CompanyDataBasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	company.AutoMigrate(&entity.Company{})

	return &DataBaseEngine{
		stockDB:   stock,
		companyDB: company,
	}, nil
}

func (engine *DataBaseEngine) Close() {
	stock, _ := engine.stockDB.DB()
	stock.Close()

	company, _ := engine.companyDB.DB()
	company.Close()
}
