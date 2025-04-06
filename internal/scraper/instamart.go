package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Vasu1712/qwiksift/server/internal/models"
	"github.com/gocolly/colly/v2"
)

type InstamartScraper struct {
    collector *colly.Collector
}

func NewInstamartScraper() *InstamartScraper {
    // In each scraper's New() function
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"))
    
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Parallelism: 1,
        Delay:       5 * time.Second,
    })

    return &InstamartScraper{
        collector: c,
    }
}

func (s *InstamartScraper) Scrape() ([]models.Product, error) {
    var products []models.Product
    baseURL := "https://www.swiggy.com"

    s.collector.OnHTML("div.product-item", func(e *colly.HTMLElement) {
        priceStr := strings.TrimSpace(e.ChildText("div.price"))
        price, _ := strconv.ParseFloat(strings.ReplaceAll(priceStr, "₹", ""), 64)

        product := models.Product{
            Name:     strings.TrimSpace(e.ChildText("h3.name")),
            Price:    price,
            ImageURL: e.ChildAttr("img", "src"),
            Provider: "Instamart",
        }
        
        products = append(products, product)
    })

	s.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://www.google.com/")
	})

    s.collector.OnError(func(r *colly.Response, err error) {
        fmt.Printf("Instamart scraper error: %v\n", err)
    })

    err := s.collector.Visit(baseURL + "/instamart")
    if err != nil {
        return nil, err
    }

    return products, nil
}

func (s *InstamartScraper) Platform() string {
    return "instamart"
}
