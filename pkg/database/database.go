package database

import (
	"chip-distribution-go/pkg/entity"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DataBaseDir  = "./DB"
	DataBasePath = "./DB/stock.db"
)

type DataBaseEngine struct {
	database *gorm.DB
}

func init() {
	_, err := os.Stat(DataBaseDir)
	if os.IsNotExist(err) {
		os.Mkdir(DataBaseDir, os.ModePerm)
	}
}

func NewDataBaseEngine() (*DataBaseEngine, error) {
	db, err := gorm.Open(sqlite.Open(DataBasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.StockInfo{})
	db.AutoMigrate(&entity.Company{})

	return &DataBaseEngine{
		database: db,
	}, nil
}

func (engine *DataBaseEngine) InsertCompanys(companyMap map[string]struct{}) {
	if companyMap == nil {
		return
	}

	for key := range companyMap {
		result := engine.database.Where("name = ?", key).First(&companyMap)
		if result.RowsAffected == 0 {
			engine.database.Create(&entity.Company{Name: key})
		}
	}
}

func (engine *DataBaseEngine) DeleteByCompanyName(companyMap map[string]struct{}) {
	if companyMap == nil {
		return
	}

	for key := range companyMap {
		engine.database.Where("name = ?", key).Delete(&entity.StockInfo{})
	}
}

func (engine *DataBaseEngine) InsertStockInfoList(infoList []entity.StockInfo) error {
	if len(infoList) > 0 {
		tx := engine.database.CreateInBatches(&infoList, 100)
		return tx.Error
	}

	return nil
}

func (engine *DataBaseEngine) Close() {
	db, _ := engine.database.DB()
	db.Close()
}