package models

type ProductResponse struct {
	Products []Product
	Total    int
	Skip     int
	Limit    int
}
