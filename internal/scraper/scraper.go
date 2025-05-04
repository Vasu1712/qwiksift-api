package scraper

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Vasu1712/qwiksift/server/internal/models"
)

type Scraper interface {
    Scrape(categoryPairs map[string]string, lat, lon string) ([]models.Product, error)
    Platform() string
}

// Define category pairs for Blinkit
var blinkitCategoryMappings = map[string]struct {
    L0Cat string
    L1Cat string
}{
    // Dairy products
    "milk":      {"14", "922"},
    "curd":      {"14", "949"},
    "bread":     {"14", "953"},
    "pav":       {"14", "953"},
    "eggs":      {"14", "1200"},
    "paneer":    {"14", "1005"},
    "cheese":    {"14", "951"},
    
    // Fruits & Vegetables
    "apple":     {"1487", "1489"},
    "banana":    {"1487", "1489"},
    "orange":    {"1487", "1489"},
    "tomato":    {"1487", "1503"},
    "potato":    {"1487", "1503"},
    "onion":     {"1487", "1503"},
    "vegetables":{"1487", "1503"},
    "fruits":    {"1487", "1489"},
    
    // Snacks
    "chips":     {"1237", "940"},
    "biscuit":   {"6", "398"},
    "chocolate": {"237", "1216"},
    
    // Default values
    "default":   {"14", "922"},
}

// Add this function to find the best match for a search query
func FindCategoryForSearchTerm(query string) (string, string) {
    query = strings.ToLower(strings.TrimSpace(query))
    
    // First try exact match
    if mapping, exists := blinkitCategoryMappings[query]; exists {
        return mapping.L0Cat, mapping.L1Cat
    }
    
    // Try partial match
    for term, mapping := range blinkitCategoryMappings {
        if strings.Contains(query, term) || strings.Contains(term, query) {
            return mapping.L0Cat, mapping.L1Cat
        }
    }
    
    // Return default if no match
    return blinkitCategoryMappings["default"].L0Cat, blinkitCategoryMappings["default"].L1Cat
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

var blinkitCategories = map[string]string{
    "14":   "922",   // Milk
    "1487": "1489", // Fresh Vegetables
    "1237": "940",  // Chips & Crisps
}

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
