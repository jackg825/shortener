# Short URL Service

This project implements a simple short URL service using Go, Gin for the HTTP framework, and PostgreSQL as the database for storing URL mappings. It demonstrates basic CRUD operations, URL shortening and expanding, and simple analytics for tracking URL visits.

## Features

- Shorten URLs with a simple API call.
- Redirect to original URLs using the shortened version.
- Basic analytics to track URL visits.
- Collision handling for shortened URL identifiers.
- In-memory caching for improved performance.

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- PostgreSQL
- Docker (optional, for running PostgreSQL in a container)
- Redis (optional, for caching)

### Installing

First, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/short-url-service.git
cd short-url-service
```

### Configuration
Create a .env file or export the following environment variables for the application configuration:

- DATABASE_URL: Connection string for the PostgreSQL database.
- REDIS_URL: Connection string for Redis (if using caching).

```bash
DATABASE_URL=postgres://user:password@localhost/urlshortener
REDIS_URL=redis://localhost:6379
```

## API Endpoints

### Shorten URL

Method: POST
Endpoint: /create
Body:
```json
{
	"longUrl": "http://example.com",
	"userId": "00000000-AAAA-BBBB-CCCC-000000000000"
}
```

### Redirect to Long URL

Method: GET
Endpoint: /{shortUrl}
Redirects to the original URL associated with the shortened version.

# Contributing

Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests to us.

# License

This project is licensed under the MIT License - see the LICENSE.md file for details.
