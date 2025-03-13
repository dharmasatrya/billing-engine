# Loan Billing System

A backend system for managing loan billings, schedules, and delinquency tracking.

## Overview

The Loan Billing System is a Go-based REST API service that handles loan management, payment processing, and delinquency tracking. It provides a robust solution for financial institutions to manage their loan portfolios.

## Features

- Loan creation and management
- 50-week payment schedule generation
- Outstanding balance tracking
- Payment processing
- Automatic delinquency detection (2+ weeks missed payments)
- Borrower management and delinquency status tracking

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Echo
- **Database**: PostgreSQL
- **ORM**: GORM
- **Validation**: go-playground/validator
- **UUID**: Google UUID
- **Documentation**: Swagger/OpenAPI

## Architecture

The system follows a clean architecture pattern with:

- **Models**: Database entities and business rules
- **Repositories**: Data access layer
- **Services**: Business logic layer
- **Handlers**: HTTP API layer
- **Scheduler**: Background processes for delinquency checks ran daily

## API Endpoints

### Borrowers
- `POST /api/borrowers`: Create a new borrower
- `GET /api/borrowers`: List all borrowers
- `GET /api/borrowers/:id`: Get borrower details
- `GET /api/borrowers/delinquent`: List delinquent borrowers

### Loans
- `POST /api/loans`: Create a new loan
- `GET /api/loans/:id`: Get loan details
- `GET /api/loans/:id/outstanding`: Get outstanding balance
- `GET /api/loans/:id/delinquent`: Check if loan is delinquent
- `POST /api/loans/:id/payment`: Make a payment

## Setup

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 13 or higher

### Environment Variables

Copy `.env.example` to `.env` and update the variables:

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=loan_billing
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
```

### Database Migration

auto migrate on startup

### Running the Application

```bash
go run cmd/api/main.go
```

### Hitting endpoint via Postman

you can import swagger.yaml to generate the Postman collection.

### Running Tests

currently only covers critical functions

```bash
# Run all tests
go test ./tests/...
```

### API Documentation

Access Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

You need to have the app running to access it.

## Loan Business Rules

1. Standard loan: 50-week term, 10% annual interest rate, equal weekly payments
2. Weekly payment = (Principal + Interest) / 50 weeks
3. A borrower is delinquent if they miss 2 or more consecutive payments
4. Payments must match the exact scheduled amount
5. Payments are applied to the earliest unpaid schedule

## Improvements to do

1. Redis caching