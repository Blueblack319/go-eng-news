package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/crazybirdz/go-eng-news/tools"
)

// Variables
const URL string = "http://www.koreaherald.com"

type article struct {
	title    string
	content  string
	createAt string
	id       string
}

var category = map[string]string{
	"020101000000": "Politics",
	"020102000000": "Social Affairs",
	"020103000000": "Foreign Affairs",
	"020107000000": "Science",
	"020109000000": "Education",
	"020203000000": "Industry",
	"020206000000": "Technology",
	"020205000000": "Transport",
	"020210000000": "Retail",
	"021901000000": "Economy",
	"021902000000": "Market",
	"021903000000": "Money",
}

var categoryIds []string

// Get category
func getCategory() int {
	var category int
	fmt.Println("Input number of category..")
	fmt.Println("0 - Politics")
	fmt.Println("1 - Social Affairs")
	fmt.Println("2 - Foreign Affairs")
	fmt.Println("3 - Science")
	fmt.Println("4 - Education")
	fmt.Println("5 - Industry")
	fmt.Println("6 - Technology")
	fmt.Println("7 - Transport")
	fmt.Println("8 - Retail")
	fmt.Println("9 - Economy")
	fmt.Println("10 - Market")
	fmt.Println("11 - Money")
	fmt.Scan(&category)
	return category
}

func main() {
	// First of all, scrape articles in first page.
	// Make many NewDocuments using goroutine
	// TODO: create mainChannel to scrape articles in all category at once
	// TODO:

	// category
	for key, _ := range category {
		categoryIds = append(categoryIds, key)
	}
	categoryIdx := getCategory()
	id := categoryIds[categoryIdx]
	articles := scrapeArticles(id)
	writeArticles(category[id], articles)
}

func getDocument(baseURL string) (*goquery.Document, error) {
	res, err := http.Get(baseURL)
	tools.CheckError(err)
	tools.CheckStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, err
}

func scrapeArticles(categoryNum string) []article {
	var articles []article
	ch := make(chan article)

	baseURL := URL + "/list.php?ct=" + categoryNum
	doc, err := getDocument(baseURL)
	tools.CheckError(err)

	articleList := doc.Find("ul.main_sec_li_only li")

	articleList.Each(func(idx int, s *goquery.Selection) {
		subURL, _ := s.Find("a").Attr("href")
		fmt.Println("Extract article:", idx)
		go extractArticle(subURL, ch)
	})

	for i := 0; i < articleList.Length(); i++ {
		extracted := <-ch
		articles = append(articles, extracted)
	}

	return articles
}

func extractArticle(subURL string, ch chan<- article) {
	baseURL := URL + subURL
	doc, err := getDocument(baseURL)
	tools.CheckError(err)

	view := doc.Find(".view")
	title := tools.CleanString(view.Find("h1.view_tit").Text())
	createAt := tools.CleanString(view.Find(".view_tit_byline_r").Text())
	content := tools.CleanString(view.Find(".view_con_t").Text())
	id := tools.GetArticleId(subURL)

	ch <- article{title: title, createAt: createAt, content: content, id: id}
}

func writeArticles(category string, articles []article) {
	length := len(articles)
	path := "articles/" + time.Now().Format("2006-January-02") + "_" + category + strconv.Itoa(length) + ".hwp"
	file, err := os.Create(path)
	tools.CheckError(err)
	defer file.Close()

	for idx, _ := range articles {
		fmt.Print(idx)
	}
}
