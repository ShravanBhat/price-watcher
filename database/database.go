package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Platform    string    `json:"platform"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PriceHistory struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

type Alert struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	OldPrice  float64   `json:"old_price"`
	NewPrice  float64   `json:"new_price"`
	Currency  string    `json:"currency"`
	Message   string    `json:"message"`
	SentAt    time.Time `json:"sent_at"`
}

func NewConnection(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize tables
	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return &DB{db}, nil
}

func initTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(500) NOT NULL,
			url TEXT NOT NULL UNIQUE,
			platform VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS price_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			price DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) DEFAULT 'INR',
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			old_price DECIMAL(10,2) NOT NULL,
			new_price DECIMAL(10,2) NOT NULL,
			currency VARCHAR(3) DEFAULT 'INR',
			message TEXT NOT NULL,
			sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_product_timestamp ON price_history(product_id, timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_timestamp ON price_history(timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_products_platform ON products(platform)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

func (db *DB) CreateProduct(name, url, platform string) (*Product, error) {
	query := `
		INSERT INTO products (name, url, platform)
		VALUES ($1, $2, $3)
		RETURNING id, name, url, platform, created_at, updated_at
	`
	
	var product Product
	err := db.QueryRow(query, name, url, platform).Scan(
		&product.ID, &product.Name, &product.URL, &product.Platform, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	return &product, nil
}

func (db *DB) GetProducts() ([]Product, error) {
	query := `SELECT id, name, url, platform, created_at, updated_at FROM products ORDER BY created_at DESC`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()
	
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.URL, &product.Platform, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	
	return products, nil
}

func (db *DB) AddPriceHistory(productID string, price float64, currency string) error {
	query := `INSERT INTO price_history (product_id, price, currency) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, productID, price, currency)
	return err
}

func (db *DB) GetLowestPriceInPeriod(productID string, days int) (float64, error) {
	query := `
		SELECT MIN(price) 
		FROM price_history 
		WHERE product_id = $1 AND timestamp >= NOW() - INTERVAL '1 day' * $2
	`
	
	var lowestPrice sql.NullFloat64
	err := db.QueryRow(query, productID, days).Scan(&lowestPrice)
	if err != nil {
		return 0, fmt.Errorf("failed to get lowest price: %w", err)
	}
	
	if !lowestPrice.Valid {
		return 0, nil
	}
	
	return lowestPrice.Float64, nil
}

func (db *DB) GetLatestPrice(productID string) (float64, error) {
	query := `
		SELECT price 
		FROM price_history 
		WHERE product_id = $1 
		ORDER BY timestamp DESC 
		LIMIT 1
	`
	
	var price sql.NullFloat64
	err := db.QueryRow(query, productID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest price: %w", err)
	}
	
	if !price.Valid {
		return 0, nil
	}
	
	return price.Float64, nil
}

func (db *DB) CreateAlert(productID string, oldPrice, newPrice float64, currency, message string) error {
	query := `INSERT INTO alerts (product_id, old_price, new_price, currency, message) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, productID, oldPrice, newPrice, currency, message)
	return err
}

func (db *DB) GetProductsByPlatform(platform string) ([]Product, error) {
	query := `SELECT id, name, url, platform, created_at, updated_at FROM products WHERE platform = $1`
	
	rows, err := db.Query(query, platform)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by platform: %w", err)
	}
	defer rows.Close()
	
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.URL, &product.Platform, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}
	
	return products, nil
}
