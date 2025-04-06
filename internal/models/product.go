package models

type Product struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Brand    string  `json:"brand"`
    Price    float64 `json:"price"`
    OldPrice float64 `json:"old_price"`
    ImageURL string  `json:"image_url"`
    Unit     string  `json:"unit"`
    Category string  `json:"category"`
    Provider string  `json:"provider"`
    Platform string  `json:"platform"`
    CategoryL0 string `json:"category_l0"`
    CategoryL1 string `json:"category_l1"`
    Slug     string  `json:"slug"`
}

