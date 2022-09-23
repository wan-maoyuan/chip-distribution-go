package entity

type ErrResponse struct {
	Message string `json:"message"`
}

type GetAllCompaniesResponse struct {
	Companies []Company `json:"companies"`
}

type GetAllStockInfosResponse struct {
	StockInfos []StockInfo `json:"stock_infos"`
}
