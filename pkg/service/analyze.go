package service

import (
	"chip-distribution-go/pkg/entity"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	AverageDistributed DistributedType = iota
	TriangularDistributed
)

type DistributedType int

func (service *ChipService) calculateConcentration(companyMap map[string]struct{}, ty DistributedType) {
	if companyMap == nil {
		return
	}

	for key := range companyMap {
		go service.checkDistributedType(key, ty)
	}
}

func (service *ChipService) checkDistributedType(name string, ty DistributedType) {
	switch ty {
	case AverageDistributed:
		service.averageDistribute(name)
	case TriangularDistributed:
		fmt.Println("TriangularDistributed")
	default:
		logrus.Errorf("DistributedType is invalid, type: %s", ty)
	}
}

type CostDistribution struct {
	DateList         []time.Time
	Chip             map[float64]float64
	ChipList         map[time.Time]map[float64]float64
	ConcentrationMap map[time.Time]*entity.StockInfo
}

// 平均分布计算
func (service *ChipService) averageDistribute(name string) {
	stockInfoList, err := service.engine.QueryStockInfoByCompanyName(name)
	if err != nil {
		logrus.Errorf("database QueryStockInfoByCompanyName function error: %v", err)
		return
	}

	cost := CostDistribution{
		DateList:         make([]time.Time, 0, len(stockInfoList)),
		Chip:             make(map[float64]float64),
		ChipList:         make(map[time.Time]map[float64]float64),
		ConcentrationMap: make(map[time.Time]*entity.StockInfo),
	}

	for index := range stockInfoList {
		cost.DateList = append(cost.DateList, stockInfoList[index].Date)
		cost.ConcentrationMap[stockInfoList[index].Date] = &stockInfoList[index]
		priceRange := []float64{}
		if stockInfoList[index].Lowest == stockInfoList[index].Highest {
			priceRange = append(priceRange, stockInfoList[index].Lowest)
		} else if stockInfoList[index].Lowest > stockInfoList[index].Highest {
			continue
		}

		tmp := stockInfoList[index].Lowest
		for tmp <= stockInfoList[index].Highest {
			priceRange = append(priceRange, tmp)
			tmp += 0.01
		}

		length := len(priceRange)
		eachVolume := stockInfoList[index].Volume / float64(length)
		for k := range cost.Chip {
			cost.Chip[k] = cost.Chip[k] * (1 - (stockInfoList[index].TurnoverRate / 100))
		}

		for _, item := range priceRange {
			_, ok := cost.Chip[item]
			if ok {
				cost.Chip[item] += eachVolume * (stockInfoList[index].TurnoverRate / 100)
			} else {
				cost.Chip[item] = eachVolume * (stockInfoList[index].TurnoverRate / 100)
			}
		}

		cost.ChipList[stockInfoList[index].Date] = deepCopyMap(cost.Chip)
	}

	// 计算成本分布
	costDistribution(&cost)

	// 将计算得到的集中度保存到数据库中
	for key := range cost.ConcentrationMap {
		if err := service.engine.UpdateStockInfoByNameDate(*cost.ConcentrationMap[key]); err != nil {
			logrus.Errorf("sqlite engine UpdateStockInfoByNameDate error: %v", err)
		}
	}
}

// 成本分布，集中度计算
func costDistribution(cost *CostDistribution) {
	for _, date := range cost.DateList {
		chip := cost.ChipList[date]
		chipKey := getKeyAndSort(chip)
		var sum float64 = 0
		var percent float64 = 0

		for i := range chip {
			sum += chip[i]
		}

		var left float64
		var right float64
		for _, j := range chipKey {
			percent += chip[j] / sum
			if percent < 0.05 {
				left = j
			} else if percent > 0.95 {
				right = j
				break
			}
		}

		if info, ok := cost.ConcentrationMap[date]; ok {
			con := (right - left) / (right + left)
			info.Concentration, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", con), 64)
		}
	}
}

func deepCopyMap(value map[float64]float64) map[float64]float64 {
	newMap := make(map[float64]float64)
	for k, v := range value {
		newMap[k] = v
	}

	return newMap
}

func getKeyAndSort(chipMap map[float64]float64) []float64 {
	keys := make([]float64, 0, len(chipMap))

	for k := range chipMap {
		keys = append(keys, k)
	}

	sort.Float64s(keys)
	return keys
}
