package entity

type ErrResponse struct {
	Message string
}

type GetAllCompaniesResponse struct {
	Companies []Company
}
