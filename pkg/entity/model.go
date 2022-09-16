package entity

import "time"

type StockInfo struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement"`
	Date          time.Time // 日期
	Code          string    // 股票代码
	Name          string    // 股票名称
	Open          float64   // 当天开盘价
	Close         float64   // 当天收盘价
	Highest       float64   // 当天最高价
	Lowest        float64   // 当天最低价
	Average       float64   // 当天平均价格
	QuoteChange   float64   // 涨跌幅
	Volume        float64   // 成交量
	Money         float64   // 成交金额
	TurnoverRate  float64   // 换手率
	Concentration float64   // 集中度
}

type Company struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"` // 公司名称
}
