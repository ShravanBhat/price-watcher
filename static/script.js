// Utility functions
function showNotification(message, type = 'info') {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.className = `notification ${type}`;
    notification.classList.remove('hidden');
    
    // Show notification
    setTimeout(() => {
        notification.classList.add('show');
    }, 100);
    
    // Hide notification after 5 seconds
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            notification.classList.add('hidden');
        }, 300);
    }, 5000);
}

function setButtonLoading(button, loading = true) {
    if (loading) {
        button.disabled = true;
        button.textContent = 'Loading...';
        button.classList.add('loading');
    } else {
        button.disabled = false;
        button.textContent = button.dataset.originalText || button.textContent;
        button.classList.remove('loading');
    }
}

// Store original button text
document.addEventListener('DOMContentLoaded', function() {
    const buttons = document.querySelectorAll('button');
    buttons.forEach(button => {
        button.dataset.originalText = button.textContent;
    });
});

// Form submission for adding products
const productForm = document.getElementById('productForm');
if (productForm) {
    productForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const submitBtn = this.querySelector('button[type="submit"]');
        const formData = new FormData(this);
        
        const productData = {
            name: formData.get('name'),
            url: formData.get('url')
        };
        
        // Validate URL
        if (!isValidUrl(productData.url)) {
            showNotification('Please enter a valid URL', 'error');
            return;
        }
        
        // Validate platform
        const platform = detectPlatform(productData.url);
        if (!platform) {
            showNotification('Unsupported platform. Please use Amazon, Flipkart, Blinkit, Zepto, Instamart, or Desidime.', 'error');
            return;
        }
        
        setButtonLoading(submitBtn, true);
        
        try {
            const response = await fetch('/api/products', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(productData)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                showNotification('Product added successfully!', 'success');
                this.reset();
                
                // Redirect to products page after a short delay
                setTimeout(() => {
                    window.location.href = '/products';
                }, 1500);
            } else {
                showNotification(result.error || 'Failed to add product', 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showNotification('Network error. Please try again.', 'error');
        } finally {
            setButtonLoading(submitBtn, false);
        }
    });
}

// Manual price scraping
async function scrapeProduct(productId) {
    const button = event.target;
    setButtonLoading(button, true);
    
    try {
        const response = await fetch(`/api/products/${productId}/scrape`, {
            method: 'POST'
        });
        
        const result = await response.json();
        
        if (response.ok) {
            showNotification(`Price scraped successfully! Current price: ₹${result.price}`, 'success');
        } else {
            showNotification(result.error || 'Failed to scrape price', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showNotification('Network error. Please try again.', 'error');
    } finally {
        setButtonLoading(button, false);
    }
}

// Delete product (placeholder - not implemented in backend yet)
async function deleteProduct(productId) {
    if (!confirm('Are you sure you want to delete this product?')) {
        return;
    }
    
    const button = event.target;
    setButtonLoading(button, true);
    
    try {
        const response = await fetch(`/api/products/${productId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            showNotification('Product deleted successfully!', 'success');
            // Remove the product card from DOM
            const productCard = button.closest('.product-card');
            productCard.remove();
        } else {
            const result = await response.json();
            showNotification(result.error || 'Failed to delete product', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showNotification('Network error. Please try again.', 'error');
    } finally {
        setButtonLoading(button, false);
    }
}

// Refresh all products
const refreshBtn = document.getElementById('refreshBtn');
if (refreshBtn) {
    refreshBtn.addEventListener('click', async function() {
        setButtonLoading(this, true);
        
        try {
            // Reload the page to refresh products
            window.location.reload();
        } catch (error) {
            console.error('Error:', error);
            showNotification('Failed to refresh products', 'error');
            setButtonLoading(this, false);
        }
    });
}

// Utility functions
function isValidUrl(string) {
    try {
        new URL(string);
        return true;
    } catch (_) {
        return false;
    }
}

function detectPlatform(url) {
    url = url.toLowerCase();
    
    if (url.includes('amazon')) return 'amazon';
    if (url.includes('flipkart')) return 'flipkart';
    if (url.includes('blinkit')) return 'blinkit';
    if (url.includes('zepto')) return 'zepto';
    if (url.includes('instamart')) return 'instamart';
    if (url.includes('desidime')) return 'desidime';
    
    return null;
}

// Add loading styles
const style = document.createElement('style');
style.textContent = `
    .btn.loading {
        opacity: 0.7;
        cursor: not-allowed;
    }
    
    .btn.loading:hover {
        transform: none !important;
        box-shadow: none !important;
    }
`;
document.head.appendChild(style);
