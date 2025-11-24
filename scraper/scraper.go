package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Scraper interface {
	ScrapePrice(url string) (float64, error)
	GetPlatformName() string
}

type BaseScraper struct {
	collector *colly.Collector
}

func NewBaseScraper() *BaseScraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.AllowURLRevisit(),
	)

	return &BaseScraper{collector: c}
}

// Amazon scraper
type AmazonScraper struct {
	*BaseScraper
}

func NewAmazonScraper() *AmazonScraper {
	return &AmazonScraper{BaseScraper: NewBaseScraper()}
}

func (a *AmazonScraper) GetPlatformName() string {
	return "amazon"
}

func (a *AmazonScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error
	var productPrice string
	var priceFound bool

	a.collector.OnHTML("#corePriceDisplay_desktop_feature_div .a-price-whole", func(e *colly.HTMLElement) {
		// Avoid overwriting if multiple similar elements are found.
		if !priceFound {
			rawPrice := e.Text

			// Use a regular expression to remove any non-digit characters (like commas).
			re := regexp.MustCompile(`[^\d]`)
			productPrice = re.ReplaceAllString(rawPrice, "")
			priceFound = true
		}
	})
	if err := a.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Amazon URL: %w", err)
	}
	price, err = strconv.ParseFloat(productPrice, 64)

	if price == 0 {
		return 0, fmt.Errorf("price not found on Amazon page")
	}

	return price, err
}

// Flipkart scraper
type FlipkartScraper struct {
	*BaseScraper
}

func NewFlipkartScraper() *FlipkartScraper {
	return &FlipkartScraper{BaseScraper: NewBaseScraper()}
}

func (f *FlipkartScraper) GetPlatformName() string {
	return "flipkart"
}

func (f *FlipkartScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error

	f.collector.OnHTML("div.Nx9bqj.CxhGGd", func(e *colly.HTMLElement) {
		priceText := strings.ReplaceAll(strings.TrimPrefix(e.Text, "₹"), ",", "")
		price, err = strconv.ParseFloat(priceText, 64)
	})

	if err := f.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Flipkart URL: %w", err)
	}

	if price == 0 {
		return 0, fmt.Errorf("price not found on Flipkart page")
	}

	return price, err
}

// Blinkit scraper
type BlinkitScraper struct {
	*BaseScraper
}

func NewBlinkitScraper() *BlinkitScraper {
	return &BlinkitScraper{BaseScraper: NewBaseScraper()}
}

func (b *BlinkitScraper) GetPlatformName() string {
	return "blinkit"
}

func (b *BlinkitScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error

	b.collector.OnHTML("span[data-testid='price']", func(e *colly.HTMLElement) {
		priceText := strings.ReplaceAll(strings.TrimPrefix(e.Text, "₹"), ",", "")
		price, err = strconv.ParseFloat(priceText, 64)
	})

	if err := b.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Blinkit URL: %w", err)
	}

	if price == 0 {
		return 0, fmt.Errorf("price not found on Blinkit page")
	}

	return price, err
}

// Zepto scraper
type ZeptoScraper struct {
	*BaseScraper
}

func NewZeptoScraper() *ZeptoScraper {
	return &ZeptoScraper{BaseScraper: NewBaseScraper()}
}

func (z *ZeptoScraper) GetPlatformName() string {
	return "zepto"
}

func (z *ZeptoScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error

	z.collector.OnHTML("span[data-testid='price']", func(e *colly.HTMLElement) {
		priceText := strings.ReplaceAll(strings.TrimPrefix(e.Text, "₹"), ",", "")
		price, err = strconv.ParseFloat(priceText, 64)
	})

	if err := z.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Zepto URL: %w", err)
	}

	if price == 0 {
		return 0, fmt.Errorf("price not found on Zepto page")
	}

	return price, err
}

// Instamart scraper
type InstamartScraper struct {
	*BaseScraper
}

func NewInstamartScraper() *InstamartScraper {
	return &InstamartScraper{BaseScraper: NewBaseScraper()}
}

func (i *InstamartScraper) GetPlatformName() string {
	return "instamart"
}

func (i *InstamartScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error

	i.collector.OnHTML("span[data-testid='price']", func(e *colly.HTMLElement) {
		priceText := strings.ReplaceAll(strings.TrimPrefix(e.Text, "₹"), ",", "")
		price, err = strconv.ParseFloat(priceText, 64)
	})

	if err := i.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Instamart URL: %w", err)
	}

	if price == 0 {
		return 0, fmt.Errorf("price not found on Instamart page")
	}

	return price, err
}

// Desidime scraper for deals
type DesidimeScraper struct {
	*BaseScraper
}

func NewDesidimeScraper() *DesidimeScraper {
	return &DesidimeScraper{BaseScraper: NewBaseScraper()}
}

func (d *DesidimeScraper) GetPlatformName() string {
	return "desidime"
}

func (d *DesidimeScraper) ScrapePrice(url string) (float64, error) {
	var price float64
	var err error

	d.collector.OnHTML("span.deal-price", func(e *colly.HTMLElement) {
		priceText := strings.ReplaceAll(strings.TrimPrefix(e.Text, "₹"), ",", "")
		price, err = strconv.ParseFloat(priceText, 64)
	})

	if err := d.collector.Visit(url); err != nil {
		return 0, fmt.Errorf("failed to visit Desidime URL: %w", err)
	}

	if price == 0 {
		return 0, fmt.Errorf("price not found on Desidime page")
	}

	return price, err
}

// ScraperFactory creates appropriate scraper based on URL
type ScraperFactory struct{}

func NewScraperFactory() *ScraperFactory {
	return &ScraperFactory{}
}

func (sf *ScraperFactory) GetScraper(url string) (Scraper, error) {
	url = strings.ToLower(url)

	switch {
	case strings.Contains(url, "amazon"):
		return NewAmazonScraper(), nil
	case strings.Contains(url, "flipkart"):
		return NewFlipkartScraper(), nil
	case strings.Contains(url, "blinkit"):
		return NewBlinkitScraper(), nil
	case strings.Contains(url, "zepto"):
		return NewZeptoScraper(), nil
	case strings.Contains(url, "instamart"):
		return NewInstamartScraper(), nil
	case strings.Contains(url, "desidime"):
		return NewDesidimeScraper(), nil
	default:
		return nil, fmt.Errorf("unsupported platform for URL: %s", url)
	}
}

// ExtractPriceFromText extracts price from text using regex
func ExtractPriceFromText(text string) (float64, error) {
	// Regex to find price patterns like ₹1,999 or 1999
	re := regexp.MustCompile(`[₹]?([0-9,]+(?:\.[0-9]{2})?)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) < 2 {
		return 0, fmt.Errorf("no price found in text: %s", text)
	}

	priceStr := strings.ReplaceAll(matches[1], ",", "")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}
