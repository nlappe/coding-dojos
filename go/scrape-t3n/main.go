package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scrape-da-t3n/models"
	"strings"
	"sync"
	"time"
)

const baseUrl = "https://t3n.de"
const overviewPagesToParse = 5
const articlesPerPage = 10

func main() {
	start := time.Now()
	log.Println("Loading News overview page")
	urlsToNewsPages := make([]string, overviewPagesToParse*articlesPerPage)

	var scrapeGroup sync.WaitGroup
	scrapeChannel := make(chan int, 10)
	for i := 0; i < overviewPagesToParse; i++ {
		i := i
		scrapeGroup.Add(1)

		scrapeChannel <- 1
		go func() {
			overviewUrl := fmt.Sprintf("%s/tag/software-entwicklung/?p=%d", baseUrl, i+1)

			res, _ := http.Get(overviewUrl)
			doc, _ := goquery.NewDocumentFromReader(res.Body)

			doc.Find(".c-pin__link").Each(func(j int, selection *goquery.Selection) {
				link, _ := selection.Attr("href")
				urlsToNewsPages[i*10+j] = link
			})
			<-scrapeChannel
			scrapeGroup.Done()
		}()
	}
	scrapeGroup.Wait()

	newsArticlesMap := map[int]models.IndexEntry{}

	_, err := os.Stat("news")
	if os.IsNotExist(err) {
		os.Mkdir("news", 0777)
	}

	var downloadWg sync.WaitGroup
	ioChannel := make(chan int, 5)

	for i, link := range urlsToNewsPages {
		downloadWg.Add(1)

		link := link
		i := i
		go func() {
			articleUrl := baseUrl + link
			res, _ := http.Get(articleUrl)
			doc, _ := goquery.NewDocumentFromReader(res.Body)

			headline := doc.Find("h2").First().Text()
			teaser := doc.Find(".u-text-teaser").First().Text()
			imgUrl, _ := doc.Find(".webfeedsFeaturedVisual").Attr("src")

			news := models.News{
				HeroImgUrl: imgUrl,
				ArticleUrl: articleUrl,
				Title:      strings.Trim(headline, "\n "),
				Teaser:     strings.Trim(teaser, "\n "),
			}

			filePath := filepath.Join("news", fmt.Sprintf("%d", i))

			indexEntry := models.IndexEntry{
				ArticleUrl: articleUrl,
				FolderPath: filePath,
			}

			_, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				os.Mkdir(filePath, 0777)
			}

			articleJson, _ := json.Marshal(news)

			os.WriteFile(filepath.Join(filePath, "data.json"), articleJson, 0777)

			if len(imgUrl) > 0 {
				ioChannel <- 1

				heroImgRes, _ := http.Get(imgUrl)
				imgBytes, _ := io.ReadAll(heroImgRes.Body)
				f, _ := os.Create(filepath.Join(filePath, "hero.jpg"))
				f.Write(imgBytes)
				f.Close()

				<-ioChannel
			}

			newsArticlesMap[i] = indexEntry
			downloadWg.Done()
		}()
	}

	downloadWg.Wait()
	close(ioChannel)

	indexJsonData, _ := json.Marshal(newsArticlesMap)
	f, _ := os.Create("news/index.json")
	defer f.Close()
	f.Write(indexJsonData)
	execTime := time.Since(start)
	log.Println("execution took ms: " + execTime.String())
}
