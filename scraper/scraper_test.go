package scraper

import (
	"testing"
)

func TestExtractPriceFromText(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPrice float64
		wantError bool
	}{
		{
			name:      "Price with rupee symbol and comma",
			input:     "₹1,999",
			wantPrice: 1999,
			wantError: false,
		},
		{
			name:      "Price without symbol",
			input:     "1999",
			wantPrice: 1999,
			wantError: false,
		},
		{
			name:      "Price with decimals",
			input:     "₹1,999.99",
			wantPrice: 1999.99,
			wantError: false,
		},
		{
			name:      "Price with multiple commas",
			input:     "₹10,99,999",
			wantPrice: 1099999,
			wantError: false,
		},
		{
			name:      "Simple integer price",
			input:     "599",
			wantPrice: 599,
			wantError: false,
		},
		{
			name:      "Price in sentence",
			input:     "The price is ₹2,499 only",
			wantPrice: 2499,
			wantError: false,
		},
		{
			name:      "No price in text",
			input:     "No price here",
			wantPrice: 0,
			wantError: true,
		},
		{
			name:      "Empty string",
			input:     "",
			wantPrice: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := ExtractPriceFromText(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("ExtractPriceFromText() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ExtractPriceFromText() unexpected error: %v", err)
				return
			}

			if price != tt.wantPrice {
				t.Errorf("ExtractPriceFromText() = %v, want %v", price, tt.wantPrice)
			}
		})
	}
}

func TestScraperFactory_GetScraper(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantPlatform string
		wantError    bool
	}{
		{
			name:         "Amazon URL",
			url:          "https://www.amazon.in/product/123",
			wantPlatform: "amazon",
			wantError:    false,
		},
		{
			name:         "Amazon URL with HTTPS",
			url:          "https://amazon.com/dp/B08N5WRWNW",
			wantPlatform: "amazon",
			wantError:    false,
		},
		{
			name:         "Flipkart URL",
			url:          "https://www.flipkart.com/product/p/itmxyz",
			wantPlatform: "flipkart",
			wantError:    false,
		},
		{
			name:         "Blinkit URL",
			url:          "https://blinkit.com/prn/product/123",
			wantPlatform: "blinkit",
			wantError:    false,
		},
		{
			name:         "Zepto URL",
			url:          "https://www.zepto.com/product/123",
			wantPlatform: "zepto",
			wantError:    false,
		},
		{
			name:         "Instamart URL",
			url:          "https://www.instamart.com/product/123",
			wantPlatform: "instamart",
			wantError:    false,
		},
		{
			name:         "Desidime URL",
			url:          "https://desidime.com/deals/product-123",
			wantPlatform: "desidime",
			wantError:    false,
		},
		{
			name:         "Case insensitive - AMAZON",
			url:          "https://WWW.AMAZON.IN/product/123",
			wantPlatform: "amazon",
			wantError:    false,
		},
		{
			name:         "Unsupported platform",
			url:          "https://www.myntra.com/product/123",
			wantPlatform: "",
			wantError:    true,
		},
		{
			name:         "Invalid URL",
			url:          "not-a-url",
			wantPlatform: "",
			wantError:    true,
		},
	}

	factory := NewScraperFactory()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scraper, err := factory.GetScraper(tt.url)

			if tt.wantError {
				if err == nil {
					t.Errorf("GetScraper() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetScraper() unexpected error: %v", err)
				return
			}

			if scraper == nil {
				t.Errorf("GetScraper() returned nil scraper")
				return
			}

			platform := scraper.GetPlatformName()
			if platform != tt.wantPlatform {
				t.Errorf("GetScraper().GetPlatformName() = %v, want %v", platform, tt.wantPlatform)
			}
		})
	}
}
