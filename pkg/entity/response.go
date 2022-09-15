package entity

type ErrResponse struct {
	Message string `json:"message"`
}

type GetAllCompaniesResponse struct {
	Companies []Company `json:"companies"`
}
