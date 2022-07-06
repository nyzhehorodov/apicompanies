package v1

type CompanyRequest struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

type CompaniesListResponse struct {
	Items []Company `json:"items"`
	Total int       `json:"total"`
}

type Company struct {
	ID      int    `json:"id"`
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}
