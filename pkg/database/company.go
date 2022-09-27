package database

import "chip-distribution-go/pkg/entity"

func (engine *DataBaseEngine) QueryCompanyByName(name string) ([]entity.Company, error) {
	defer engine.companyRW.RUnlock()
	engine.companyRW.RLock()

	var companies = []entity.Company{}
	result := engine.companyDB.Where("name = ?", name).Find(&companies)
	return companies, result.Error
}

func (engine *DataBaseEngine) QueryAllCompanies() ([]entity.Company, error) {
	defer engine.companyRW.RUnlock()
	engine.companyRW.RLock()

	var companies []entity.Company
	result := engine.companyDB.Find(&companies)
	return companies, result.Error
}

func (engine *DataBaseEngine) InsertCompanies(companyMap map[string]struct{}) {
	if companyMap == nil {
		return
	}

	for key := range companyMap {
		companies, err := engine.QueryCompanyByName(key)
		if err == nil && len(companies) == 0 {
			result := engine.companyDB.Where("name = ?", key).First(&entity.Company{})
			if result.RowsAffected == 0 {
				engine.InsertCompany(key)
			}
		}
	}
}

func (engine *DataBaseEngine) InsertCompany(companyName string) {
	defer engine.companyRW.Unlock()
	engine.companyRW.Lock()

	engine.companyDB.Create(&entity.Company{Name: companyName})
}
