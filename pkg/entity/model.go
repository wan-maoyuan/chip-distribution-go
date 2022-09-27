package entity

import "time"

type StockInfo struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Date          time.Time `json:"date"`              // 日期
	Code          string    `json:"code"`              // 股票代码
	Name          string    `gorm:"index" json:"name"` // 股票名称
	Open          float64   `json:"open"`              // 当天开盘价
	Close         float64   `json:"close"`             // 当天收盘价
	Highest       float64   `json:"hightest"`          // 当天最高价
	Lowest        float64   `json:"lowest"`            // 当天最低价
	Average       float64   `json:"average"`           // 当天平均价格
	QuoteChange   float64   `json:"quote_change"`      // 涨跌幅
	Volume        float64   `json:"volume"`            // 成交量
	Money         float64   `json:"money"`             // 成交金额
	TurnoverRate  float64   `json:"turnover_rate"`     // 换手率
	Concentration float64   `json:"concentration"`     // 集中度
}

type Company struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"` // 公司名称
}
