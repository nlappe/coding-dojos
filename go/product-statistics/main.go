package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"product-statistics/models"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	defer func() {
		execTime := time.Since(start)
		log.Println("execution took ms: " + execTime.String())
	}()

	args := os.Args
	if len(args) < 2 {
		log.Println("No Product name given.")
		return
	}

	productName := args[1]

	res, _ := os.ReadFile("products.json")

	var productResponse models.ProductResponse
	json.Unmarshal(res, &productResponse)

	matchedProduct := findProduct(productResponse.Products, productName)

	if matchedProduct == nil {
		log.Println(fmt.Sprintf("Product with name '%s' was not found in the dataset. The following Products are available:", productName))
		productNames := getProductNames(productResponse.Products)

		log.Println(strings.Join(productNames, ", "))
		return
	}

	productStats := models.ProductStats{
		Name:            matchedProduct.Title,
		Price:           fmt.Sprintf("%d€", matchedProduct.Price),
		DiscountedPrice: fmt.Sprintf("%.2f€", float32(matchedProduct.Price)*(1.0-(matchedProduct.DiscountPercentage/100))),
		DiscountPercent: matchedProduct.DiscountPercentage,
		Rating:          strings.Repeat("⭐ ", int(matchedProduct.Rating)),
		Category:        matchedProduct.Category,
		Brand:           matchedProduct.Brand,
		Stock:           matchedProduct.Stock,
	}

	log.Println(fmt.Sprintf("\n\n"+
		"Produkt: %s \n"+
		"Preis: %s \n"+
		"Reduzierter Preis: %s \n"+
		"Bewertung: %s \n"+
		"Kategorie: %s \n"+
		"Marke: %s \n"+
		"", productStats.Name, productStats.Price, productStats.DiscountedPrice, productStats.Rating, productStats.Category, productStats.Brand))
}

func findProduct(productArray []models.Product, productName string) *models.Product {
	for _, product := range productArray {
		if strings.ToLower(product.Title) == strings.ToLower(productName) {
			return &product
		}
	}

	return nil
}

func getProductNames(products []models.Product) []string {
	names := make([]string, len(products))

	for i, product := range products {
		names[i] = product.Title
	}

	return names
}
