# Price Watcher

A comprehensive price monitoring application built in Go that tracks product prices across multiple e-commerce platforms and sends alerts via Telegram when prices drop.

## Features

- **Multi-Platform Support**: Monitor prices on Amazon and Flipkart
- **Automated Scraping**: Scheduled price scraping with configurable intervals
- **Smart Alerts**: Telegram notifications when prices drop to their lowest in the configured period
- **Price History**: Track price changes over time (configurable, default: 30 days)
- **Web Interface**: Clean, responsive web UI for managing products
- **Concurrent Processing**: Worker pool architecture for efficient scraping
- **Database Storage**: PostgreSQL backend for reliable data persistence


## Technology Stack

- **Backend**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Web Scraping**: Colly
- **Scheduling**: Cron
- **Telegram Integration**: go-telegram-bot-api
- **Frontend**: HTML5, CSS3, JavaScript (ES6+)
- **Styling**: Custom CSS with responsive design

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Telegram Bot Token (optional, for alerts)

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ShravanBhat/price-watcher
   cd price-watcher
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE price_watcher;
   CREATE USER price_watcher_user WITH PASSWORD 'your_password';
   GRANT ALL PRIVILEGES ON DATABASE price_watcher TO price_watcher_user;
   ```

4. **Configure environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost:5432/price_watcher?sslmode=disable` |
| `TELEGRAM_TOKEN` | Telegram bot token | (empty) |
| `TELEGRAM_CHAT_ID` | Telegram chat ID for alerts | (empty) |
| `SERVER_PORT` | HTTP server port | `8080` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout (seconds) | `30` |
| `SCRAPING_INTERVAL` | Price scraping interval (seconds) | `3600` (1 hour) |
| `PRICE_HISTORY_DAYS` | Days to keep price history | `30` |

### Telegram Bot Setup

1. Create a new bot via [@BotFather](https://t.me/botfather)
2. Get the bot token
3. Start a chat with your bot
4. Get your chat ID (you can use [@userinfobot](https://t.me/userinfobot))
5. Add both values to your `.env` file

## Usage

### Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Add products by providing:
   - Product name
   - Product URL from supported platforms
3. View all products at `/products`
4. Manually trigger price scraping for individual products

### API Endpoints

- `POST /api/products` - Add a new product
- `GET /api/products` - Get all products
- `DELETE /api/products/:id` - Delete a product
- `POST /api/products/:id/scrape` - Manually scrape price


### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o price-watcher main.go
```

## Alert Logic

The application sends alerts when:

1. **Price Change Detected**: Current price differs from previous price
2. **New Low Price**: Current price is the lowest in the configured period (default: 30 days)
3. **No Duplicate Alerts**: Alerts are only sent for actual price changes

Alert messages include:
- Product name and platform
- Previous and current prices
- Amount saved
- Whether it's the lowest price in the period
- Direct link to the product

## Database Schema

### Tables

- **`products`**: Product information and metadata
- **`price_history`**: Historical price data
- **`alerts`**: Sent alert records

## Future Enhancements

- [ ] Email notifications
- [ ] Price trend analysis
- [ ] Multiple user support
- [ ] Mobile app
- [ ] Advanced filtering and search
- [ ] Price comparison across platforms
- [ ] Historical price charts
- [ ] Export functionality
