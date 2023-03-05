package ui_renderer

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"image/color"
	"log"
	"product-statistics/models"
	"strings"
)

func ShowProductStats(product models.ProductStats) {
	uiApp := app.New()
	window := uiApp.NewWindow("Product-Stats")
	gridLayout := layout.NewGridLayout(2)

	log.Println(strings.ReplaceAll(product.Rating, " ", ""))

	gridContainer := container.New(gridLayout,
		canvas.NewText("Product", color.White),
		canvas.NewText(product.Name, color.White),
		canvas.NewText("Price", color.White),
		canvas.NewText(product.Price, color.White),
		canvas.NewText("Special Offer", color.White),
		canvas.NewText(product.DiscountedPrice, color.White),
		canvas.NewText("Discount", color.White),
		canvas.NewText(fmt.Sprintf("%.2f", product.DiscountPercent)+"%", color.White),
		canvas.NewText("Rating", color.White),
		canvas.NewText(product.Rating, color.White),
		canvas.NewText("Category", color.White),
		canvas.NewText(product.Category, color.White),
		canvas.NewText("Brand", color.White),
		canvas.NewText(product.Brand, color.White),
		canvas.NewText("Stock", color.White),
		canvas.NewText(fmt.Sprintf("%d", product.Stock), color.White),
	)

	window.SetContent(gridContainer)

	window.ShowAndRun()
}
