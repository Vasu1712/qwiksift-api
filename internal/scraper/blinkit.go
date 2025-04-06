package scraper

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strconv"
    "time"
    
    "github.com/Vasu1712/qwiksift/server/internal/models"
)

// BlinkitScraper handles API requests to Blinkit
type BlinkitScraper struct {
    baseURL string
    client  *http.Client
}

// NewBlinkitScraper creates a new Blinkit scraper instance
func NewBlinkitScraper() *BlinkitScraper {
    return &BlinkitScraper{
        baseURL: "https://blinkit.com",
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// FetchCategoryProducts retrieves products from a specific category and returns the raw JSON
func (s *BlinkitScraper) FetchCategoryProducts(l0Cat, l1Cat, lat, lon string) ([]byte, error) {
    // Construct the URL with the category parameters
    url := fmt.Sprintf("%s/v2/listing?l0_cat=%s&l1_cat=%s",
        s.baseURL, l0Cat, l1Cat)
    
    // Create a new request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }
    
    // Add required headers - match exactly what's in your curl command
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Lat", lat)
    req.Header.Set("Lon", lon)
    
    // Make the request
    log.Printf("Requesting URL: %s with headers Lat: %s, Lon: %s", url, lat, lon)
    resp, err := s.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()
    
    // Log response details
    log.Printf("Status Code for %s: %d", url, resp.StatusCode)
    
    // Check status code
    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
    }
    
    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }
    
    return body, nil
}

// Scrape fetches products from multiple categories and returns them as a slice of models.Product
func (s *BlinkitScraper) Scrape(categoryPairs map[string]string, lat, lon string) ([]models.Product, error) {
    var allProducts []models.Product
    
    for l0CatStr, l1CatStr := range categoryPairs {
        jsonData, err := s.FetchCategoryProducts(l0CatStr, l1CatStr, lat, lon)
        if err != nil {
            fmt.Printf("Error fetching category %s_%s: %v\n", l0CatStr, l1CatStr, err)
            continue
        }
        
        // For simplicity, let's pass through the raw JSON response
        var response map[string]interface{}
        if err := json.Unmarshal(jsonData, &response); err != nil {
            fmt.Printf("Error unmarshaling JSON: %v\n", err)
            continue
        }
        
        // Extract products from the response
        if data, ok := response["data"].(map[string]interface{}); ok {
            if productsData, ok := data["products"].([]interface{}); ok {
                for _, productData := range productsData {
                    if prod, ok := productData.(map[string]interface{}); ok {
                        product := models.Product{
                            Platform:   s.Platform(),
                            CategoryL0: l0CatStr,
                            CategoryL1: l1CatStr,
                            Provider:   "Blinkit",
                        }
                        
                        // Extract name
                        if name, ok := prod["name"].(string); ok {
                            product.Name = name
                        }
                        
                        // Extract price
                        if price, ok := prod["price"].(float64); ok {
                            product.Price = price
                        }
                        
                        // Extract image URL
                        if images, ok := prod["images"].([]interface{}); ok && len(images) > 0 {
                            if img, ok := images[0].(map[string]interface{}); ok {
                                if url, ok := img["url"].(string); ok {
                                    product.ImageURL = url
                                }
                            }
                        }
                        
                        // Extract ID
                        if id, ok := prod["id"].(float64); ok {
                            product.ID = strconv.FormatFloat(id, 'f', 0, 64)
                        }
                        
                        allProducts = append(allProducts, product)
                    }
                }
            }
        }
    }
    
    return allProducts, nil
}

// Platform returns the platform name
func (s *BlinkitScraper) Platform() string {
    return "blinkit"
}
