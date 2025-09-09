package scheduler

import (
	"fmt"
	"log"
	"sync"

	"price-watcher/config"
	"price-watcher/database"
	"price-watcher/scraper"
	"price-watcher/telegram"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron     *cron.Cron
	db       *database.DB
	tgBot    *telegram.Bot
	config   *config.Config
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewScheduler(db *database.DB, tgBot *telegram.Bot, cfg *config.Config) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		db:       db,
		tgBot:    tgBot,
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	log.Println("Starting price watcher scheduler...")

	// Schedule price scraping every hour (or as configured)
	s.cron.AddFunc(fmt.Sprintf("0 */%d * * * *", int(s.config.ScrapingInterval.Hours())), s.scrapeAllProducts)

	// Start the cron scheduler
	s.cron.Start()

	// Run initial scraping
	go s.scrapeAllProducts()
}

func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	close(s.stopChan)
	s.cron.Stop()
	s.wg.Wait()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) scrapeAllProducts() {
	log.Println("Starting scheduled price scraping...")

	// Get all products from database
	products, err := s.db.GetProducts()
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		return
	}

	if len(products) == 0 {
		log.Println("No products to scrape")
		return
	}

	// Create worker pool for concurrent scraping
	workerCount := 5
	productChan := make(chan database.Product, len(products))

	// Start workers
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.priceScrapingWorker(productChan)
	}

	// Send products to workers
	for _, product := range products {
		select {
		case productChan <- product:
		case <-s.stopChan:
			close(productChan)
			return
		}
	}

	close(productChan)
	s.wg.Wait()

	log.Println("Price scraping completed")
}

func (s *Scheduler) priceScrapingWorker(productChan <-chan database.Product) {
	defer s.wg.Done()

	for product := range productChan {
		select {
		case <-s.stopChan:
			return
		default:
			s.scrapeProductPrice(product)
		}
	}
}

func (s *Scheduler) scrapeProductPrice(product database.Product) {
	log.Printf("Scraping price for product: %s (%s)", product.Name, product.Platform)

	// Get appropriate scraper for the platform
	factory := scraper.NewScraperFactory()
	scraper, err := factory.GetScraper(product.URL)
	if err != nil {
		log.Printf("Failed to get scraper for %s: %v", product.URL, err)
		return
	}

	// Scrape current price
	currentPrice, err := scraper.ScrapePrice(product.URL)
	if err != nil {
		log.Printf("Failed to scrape price for %s: %v", product.URL, err)
		return
	}

	// Check if we should send an alert
	if err := s.checkAndSendAlert(product, currentPrice); err != nil {
		log.Printf("Failed to check/send alert for %s: %v", product.ID, err)
	}

	// Add price to history
	if err := s.db.AddPriceHistory(product.ID, currentPrice, "INR"); err != nil {
		log.Printf("Failed to add price history for %s: %v", product.ID, err)
		return
	}

	log.Printf("Successfully scraped price for %s: ₹%.2f", product.Name, currentPrice)
}

func (s *Scheduler) checkAndSendAlert(product database.Product, currentPrice float64) error {
	// Get the previous price
	previousPrice, err := s.db.GetLatestPrice(product.ID)
	if err != nil {
		// If no previous price, this is the first scrape
		return nil
	}

	// If prices are the same, no need to send alert
	if currentPrice == previousPrice {
		return nil
	}

	// Get the lowest price in the configured period
	lowestPrice, err := s.db.GetLowestPriceInPeriod(product.ID, s.config.PriceHistoryDays)
	if err != nil {
		log.Printf("Failed to get lowest price for %s: %v", product.ID, err)
		return err
	}

	// Check if current price is the lowest in the period
	if lowestPrice == 0 || currentPrice <= lowestPrice {
		// Send alert
		message := fmt.Sprintf(
			"🚨 PRICE DROP ALERT! 🚨\n\n"+
				"Product: %s\n"+
				"Platform: %s\n"+
				"Previous Price: ₹%.2f\n"+
				"Current Price: ₹%.2f\n"+
				"Savings: ₹%.2f\n"+
				"Lowest in %d days: %s\n\n"+
				"🔗 %s",
			product.Name,
			product.Platform,
			previousPrice,
			currentPrice,
			previousPrice-currentPrice,
			s.config.PriceHistoryDays,
			func() string {
				if currentPrice == lowestPrice {
					return "YES! 🎉"
				}
				return fmt.Sprintf("₹%.2f", lowestPrice)
			}(),
			product.URL,
		)

		// Send to Telegram
		if err := s.tgBot.SendMessage(message); err != nil {
			log.Printf("Failed to send Telegram alert: %v", err)
			return err
		}

		// Store alert in database
		if err := s.db.CreateAlert(product.ID, previousPrice, currentPrice, "INR", message); err != nil {
			log.Printf("Failed to store alert: %v", err)
		}

		log.Printf("Alert sent for %s: Price dropped from ₹%.2f to ₹%.2f",
			product.Name, previousPrice, currentPrice)
	}

	return nil
}

// ManualScrape allows manual triggering of price scraping for a specific product
func (s *Scheduler) ManualScrape(productID string) error {
	// Get product details
	products, err := s.db.GetProducts()
	if err != nil {
		return fmt.Errorf("failed to get products: %w", err)
	}

	var targetProduct database.Product
	found := false
	for _, product := range products {
		if product.ID == productID {
			targetProduct = product
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("product not found: %s", productID)
	}

	// Scrape the product
	s.scrapeProductPrice(targetProduct)
	return nil
}
