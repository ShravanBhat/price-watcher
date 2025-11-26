package database

import (
	"os"
	"testing"
)

// Note: These tests require a running PostgreSQL database.
// Set TEST_DATABASE_URL environment variable to run these tests.
// Example: TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/price_watcher_test?sslmode=disable

func getTestDB(t *testing.T) *DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database tests")
	}

	db, err := NewConnection(dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func TestDeleteProduct_Success(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Create a test product
	product, err := db.CreateProduct("Test Product for Deletion", "https://www.amazon.in/test-delete", "amazon")
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	// Delete the product
	err = db.DeleteProduct(product.ID)
	if err != nil {
		t.Errorf("DeleteProduct() unexpected error: %v", err)
	}

	// Verify product was deleted by trying to get all products
	products, err := db.GetProducts()
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}

	// Check that the deleted product is not in the list
	for _, p := range products {
		if p.ID == product.ID {
			t.Errorf("Product %s still exists after deletion", product.ID)
		}
	}
}

func TestDeleteProduct_NotFound(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Try to delete a product with a non-existent UUID
	fakeUUID := "00000000-0000-0000-0000-000000000000"
	err := db.DeleteProduct(fakeUUID)

	if err == nil {
		t.Error("DeleteProduct() expected error for non-existent product, got nil")
	}

	expectedError := "product not found"
	if err != nil && !contains(err.Error(), expectedError) {
		t.Errorf("DeleteProduct() error = %v, want error containing %q", err, expectedError)
	}
}

func TestDeleteProduct_InvalidID(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Try to delete with an invalid UUID format
	invalidID := "not-a-valid-uuid"
	err := db.DeleteProduct(invalidID)

	if err == nil {
		t.Error("DeleteProduct() expected error for invalid UUID, got nil")
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Create a test product
	productName := "Integration Test Product"
	productURL := "https://www.flipkart.com/test-product"
	platform := "flipkart"

	product, err := db.CreateProduct(productName, productURL, platform)
	if err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}

	// Clean up after test
	defer db.DeleteProduct(product.ID)

	// Verify the created product
	if product.Name != productName {
		t.Errorf("CreateProduct() name = %v, want %v", product.Name, productName)
	}
	if product.URL != productURL {
		t.Errorf("CreateProduct() url = %v, want %v", product.URL, productURL)
	}
	if product.Platform != platform {
		t.Errorf("CreateProduct() platform = %v, want %v", product.Platform, platform)
	}

	// Get all products and verify our product is in the list
	products, err := db.GetProducts()
	if err != nil {
		t.Fatalf("GetProducts() error = %v", err)
	}

	found := false
	for _, p := range products {
		if p.ID == product.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Created product not found in GetProducts() result")
	}
}

func TestAddPriceHistory(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	// Create a test product
	product, err := db.CreateProduct("Price History Test Product", "https://www.amazon.in/test-price-history", "amazon")
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	defer db.DeleteProduct(product.ID)

	// Add price history with delta
	price := 1000.00
	delta := -50.00
	err = db.AddPriceHistory(product.ID, price, delta, "INR")
	if err != nil {
		t.Fatalf("AddPriceHistory() error = %v", err)
	}

	// Verify the price history was added
	// We need to query the database directly or add a GetPriceHistory method to DB for testing
	// For now, let's use GetLatestPrice to verify at least the price
	latestPrice, err := db.GetLatestPrice(product.ID)
	if err != nil {
		t.Fatalf("GetLatestPrice() error = %v", err)
	}

	if latestPrice != price {
		t.Errorf("GetLatestPrice() = %v, want %v", latestPrice, price)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
