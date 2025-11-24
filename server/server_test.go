package server

import (
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantPlatform string
	}{
		{
			name:         "Amazon URL",
			url:          "https://www.amazon.in/product/123",
			wantPlatform: "amazon",
		},
		{
			name:         "Amazon with uppercase",
			url:          "https://WWW.AMAZON.COM/dp/B08N5WRWNW",
			wantPlatform: "amazon",
		},
		{
			name:         "Flipkart URL",
			url:          "https://www.flipkart.com/product/p/itmxyz",
			wantPlatform: "flipkart",
		},
		{
			name:         "Blinkit URL",
			url:          "https://blinkit.com/prn/product/123",
			wantPlatform: "blinkit",
		},
		{
			name:         "Zepto URL",
			url:          "https://www.zepto.com/product/123",
			wantPlatform: "zepto",
		},
		{
			name:         "Instamart URL",
			url:          "https://www.instamart.com/product/123",
			wantPlatform: "instamart",
		},
		{
			name:         "Desidime URL",
			url:          "https://desidime.com/deals/product-123",
			wantPlatform: "desidime",
		},
		{
			name:         "Mixed case Flipkart",
			url:          "https://www.FlipKart.com/product",
			wantPlatform: "flipkart",
		},
		{
			name:         "Unsupported platform",
			url:          "https://www.myntra.com/product/123",
			wantPlatform: "",
		},
		{
			name:         "Random URL",
			url:          "https://www.example.com",
			wantPlatform: "",
		},
		{
			name:         "Empty URL",
			url:          "",
			wantPlatform: "",
		},
	}

	// Create a server instance to test the method
	s := &Server{}

	for _,tt:=range tests{
		t.Run(tt.name, func(t *testing.T){
			platform:=s.detectPlatform(tt.url)
			if platform!=tt.wantPlatform{
				t.Errorf("detectPlatform() = %v, want %v", platform, tt.wantPlatform)
			}
		})
	}
}
