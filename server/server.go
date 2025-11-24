package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"price-watcher/config"
	"price-watcher/database"
	"price-watcher/scraper"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	db     *database.DB
	config *config.Config
	server *http.Server
}

func NewServer(db *database.DB, cfg *config.Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		router: router,
		db:     db,
		config: cfg,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Serve static files
	s.router.Static("/static", "./static")
	s.router.LoadHTMLGlob("templates/*")

	// API routes
	api := s.router.Group("/api")
	{
		api.POST("/products", s.createProduct)
		api.GET("/products", s.getProducts)
		api.DELETE("/products/:id", s.deleteProduct)
		api.POST("/products/:id/scrape", s.manualScrape)
	}

	// Web routes
	s.router.GET("/", s.indexPage)
	s.router.GET("/products", s.productsPage)
}

func (s *Server) Start() error {
	addr := ":" + s.config.ServerPort
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	fmt.Printf("Server starting on port %s\n", s.config.ServerPort)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Price Watcher - Add Product",
	})
}

func (s *Server) productsPage(c *gin.Context) {
	products, err := s.db.GetProducts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to load products",
		})
		return
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Price Watcher - Products",
		"products": products,
	})
}

func (s *Server) createProduct(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		URL  string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine platform from URL
	platform := s.detectPlatform(req.URL)
	if platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported platform"})
		return
	}

	// Create product
	product, err := s.db.CreateProduct(req.Name, req.URL, platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (s *Server) getProducts(c *gin.Context) {
	products, err := s.db.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID required"})
		return
	}

	// Delete the product from database
	err := s.db.DeleteProduct(id)
	if err != nil {
		// Check if product was not found
		if strings.Contains(err.Error(), "product not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"id":      id,
	})
}

func (s *Server) manualScrape(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID required"})
		return
	}

	// Get product details
	products, err := s.db.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var targetProduct database.Product
	found := false
	for _, product := range products {
		if product.ID == id {
			targetProduct = product
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Scrape price
	factory := scraper.NewScraperFactory()
	scraper, err := factory.GetScraper(targetProduct.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := scraper.ScrapePrice(targetProduct.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add to price history
	if err := s.db.AddPriceHistory(targetProduct.ID, price, "INR"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Price scraped successfully",
		"price":   price,
		"product": targetProduct.Name,
	})
}

func (s *Server) detectPlatform(url string) string {
	url = strings.ToLower(url)

	switch {
	case strings.Contains(url, "amazon"):
		return "amazon"
	case strings.Contains(url, "flipkart"):
		return "flipkart"
	case strings.Contains(url, "blinkit"):
		return "blinkit"
	case strings.Contains(url, "zepto"):
		return "zepto"
	case strings.Contains(url, "instamart"):
		return "instamart"
	case strings.Contains(url, "desidime"):
		return "desidime"
	default:
		return ""
	}
}
