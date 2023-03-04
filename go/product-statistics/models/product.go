package models

type Product struct {
	Id                 uint
	Title              string
	Description        string
	Price              uint
	DiscountPercentage float32
	Rating             float32
	Stock              uint
	Brand              string
	Category           string
	Thumbnail          string
	Images             []string
}
