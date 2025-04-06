package scraper

import (
    "math/rand"
    "sync"
    "time"
    
    "github.com/Vasu1712/qwiksift/server/internal/models"
)

type Scraper interface {
    Scrape(categoryPairs map[string]string, lat, lon string) ([]models.Product, error)
    Platform() string
}

// Define category pairs for Blinkit
var blinkitCategories = map[string]string{
    "14":   "922",   // Milk
    "1487": "1489", // Fresh Vegetables
    "1237": "940",  // Chips & Crisps
}

// Define location coordinates
var locations = []struct {
    Lat string
    Lon string
}{
    {"28.7041", "77.1025"}, // Delhi
    {"19.0760", "72.8777"}, // Mumbai
    {"12.9716", "77.5946"}, // Bangalore
    {"12.9337741", "77.7006304"}, // Your specific location
}

var (
    scrapers = []Scraper{
        NewBlinkitScraper(),
        // NewZeptoScraper(),
        // NewInstamartScraper(),
    }
    
    requestDelay = time.Duration(rand.Intn(5)+3) * time.Second
)

// ScrapeAll calls Scrape() on all registered scrapers
func ScrapeAll() ([]models.Product, error) {
    var wg sync.WaitGroup
    results := make(chan []models.Product)
    errChan := make(chan error, len(scrapers))
    
    for _, s := range scrapers {
        wg.Add(1)
        go func(s Scraper) {
            defer wg.Done()
            time.Sleep(requestDelay) // Rate limiting
            
            // Always use your specific coordinates
            lat, lon := "12.9337741", "77.7006304"
            
            products, err := s.Scrape(blinkitCategories, lat, lon)
            if err != nil {
                errChan <- err
                return
            }
            results <- products
        }(s)
    }
    
    go func() {
        wg.Wait()
        close(results)
        close(errChan)
    }()
    
    var allProducts []models.Product
    for products := range results {
        allProducts = append(allProducts, products...)
    }
    
    if len(errChan) > 0 {
        return nil, <-errChan
    }
    
    return allProducts, nil
}

// GetRandomLocation returns a random location for geolocation headers
func GetRandomLocation() (string, string) {
    loc := locations[rand.Intn(len(locations))]
    return loc.Lat, loc.Lon
}
