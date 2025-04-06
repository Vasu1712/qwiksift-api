package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Vasu1712/qwiksift/server/internal/models"
	"github.com/gocolly/colly/v2"
)

type ZeptoScraper struct {
    collector *colly.Collector
}

func NewZeptoScraper() *ZeptoScraper {
    // In each scraper's New() function
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"))
    
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Parallelism: 1,
        Delay:       5 * time.Second,
    })

    return &ZeptoScraper{
        collector: c,
    }
}

func (s *ZeptoScraper) Scrape() ([]models.Product, error) {
    var products []models.Product
    baseURL := "https://www.zeptonow.com"

    s.collector.OnHTML("div.product-item", func(e *colly.HTMLElement) {
        priceStr := strings.TrimSpace(e.ChildText("div.price-box"))
        price, _ := strconv.ParseFloat(strings.ReplaceAll(priceStr, "₹", ""), 64)

        product := models.Product{
            Name:     strings.TrimSpace(e.ChildText("h2.product-name")),
            Price:    price,
            ImageURL: e.ChildAttr("img", "data-src"),
            Provider: "Zepto",
        }
        
        products = append(products, product)
    })

	s.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://www.google.com/")
	})	

    s.collector.OnError(func(r *colly.Response, err error) {
        fmt.Printf("Zepto scraper error: %v\n", err)
    })

    err := s.collector.Visit(baseURL + "/products")
    if err != nil {
        return nil, err
    }

    return products, nil
}

func (s *ZeptoScraper) Platform() string {
    return "zepto"
}
