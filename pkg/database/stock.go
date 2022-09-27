package database

import "chip-distribution-go/pkg/entity"

func (engine *DataBaseEngine) QueryStockInfoByCompanyName(companyName string) ([]entity.StockInfo, error) {
	defer engine.stockRW.RUnlock()
	engine.stockRW.RLock()

	var infos []entity.StockInfo
	result := engine.stockDB.Where("name = ?", companyName).Order("date").Find(&infos)
	return infos, result.Error
}

func (engine *DataBaseEngine) DeleteStockInfoByCompanyName(companyMap map[string]struct{}) {
	if companyMap == nil {
		return
	}

	defer engine.stockRW.Unlock()
	engine.stockRW.Lock()

	for key := range companyMap {
		engine.stockDB.Where("name = ?", key).Delete(&entity.StockInfo{})
	}
}

func (engine *DataBaseEngine) InsertStockInfoList(infoList []entity.StockInfo) error {
	defer engine.stockRW.Unlock()
	engine.stockRW.Lock()

	if len(infoList) > 0 {
		tx := engine.stockDB.CreateInBatches(&infoList, 100)
		return tx.Error
	}

	return nil
}

func (engine *DataBaseEngine) UpdateStockInfoByNameDate(info entity.StockInfo) error {
	defer engine.stockRW.Unlock()
	engine.stockRW.Lock()

	result := engine.stockDB.
		Model(&info).
		Update("concentration", info.Concentration)

	return result.Error
}
