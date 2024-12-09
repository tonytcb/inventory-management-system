# Inventory Management system

## Overview

Backend service that simulates an inventory management system for a company providing cross-border liquidity services. The company operates liquidity pools in multiple currencies and manages currency transfers for users sending funds between different currencies.

The system will support the following currencies:
- USD, EUR, JPY, GBP, and AUD.

## Objective

The key objectives are to create a system that:
- Manages liquidity pools across various currencies.
- Tracks inventory, revenue, and FX rates for each transaction.
- Dynamically rebalances liquidity pools based on transaction data to maintain optimal balances.

## Design solutions

The application architecture follows the principles of the **Clean Architecture**, originally described by Robert C. Martin. The foundation of this kind of architecture is the dependency injection, producing systems that are independent of external agents, testable and easier to maintain.
You can read more about Clean Architecture [here](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

### Ledger model

Both `currency pool` and `fx rate` operates following a ledger approach, to keep record of all changes of all time.

### Create transfer

Transfer is created as `pending` status, and will be processed asynchronously by a worker.

### Settlement transfer

Via a worker queue, the pending transfer will be updated to `completed` once the transfer is done.

Also, a transaction with margin (revenue) is stored.

### Rebalancing pools

Rebalancing runs every X time given an initial configurations.

The formula considers:
- the current balance of the pools
- transaction volume between currencies
- threshold to trigger rebalancing, given initial configuration

When **imbalance** between currencies is greater than **[volume * threshold]**, the system will rebalance the pools.

## Tools

- [Docker](https://www.docker.com/)
- [Golang 1.23](https://golang.org/)
- [Gin HTTP Router](https://github.com/gin-gonic/gin)
- Postgres

## Configurations

All next configurations are set in `.env` file.

- `DATABASE_URL`: Database url
- `REBALANCE_POOL_THRESHOLD_PERCENT`: Threshold percent to trigger rebalance
- `REBALANCE_CHECK_INTERVAL`: Time interval to check if rebalance is needed 
- `CURRENCIES_ENABLED`: List of currencies enabled in the system
- `REVENUE_MARGIN_PERCENT`: Margin percent to be added to the transfer via transaction record

## How to run

1. Create .env from .env.example
2. Run `make up` or `docker-compose up app`
3. Run `node mockFxRateSender.js http://localhost:8080/fx-rate` to start sending mock fx rates

## HTTP API

### Update FX Rate

Request
```http
POST /fx-rate

{
    "pair": "AUD/JPY",
    "rate": "82.6666",
    "timestamp": "2024-11-11T11:22:18.123Z"
}
```

Response
```http
HTTP/1.1 200 OK
```

### Create transfer

Request
```http
POST /transfer

{
    "from_account": {
        "currency": "USD"
    },
    "to_account": {
        "currency": "EUR"
    },
    "amount": "100"
}
```

Response
```http
HTTP/1.1 201 OK

{
    "transfer": {
        "ID": 1,
        "converted_amount": "108.47",
        "final_amount": "109.5547",
        "original_amount": "100",
        "from": {
            "currency": "USD"
        },
        "to": {
            "currency": "EUR"
        },
        "created_at": "2024-12-09T17:40:43.154887Z"
    }
}
```
