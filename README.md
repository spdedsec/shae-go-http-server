# SHAE Go HTTP Server

A minimal HTTP server written in Go.

## Features
- GET / — health check
- GET /api/data — return JSON data
- POST /api/echo — echo request body
- GET /items — list items
- POST /items — create item
- GET /items/:id — get item
- DELETE /items/:id — delete item

## Install
```bash
./scripts/install.sh
# Or manually:
go build -o shae-http-server src/main.go
cp shae-http-server ~/.local/bin/
```

## Usage
```bash
shae-http-server
# Server listens on http://0.0.0.0:8080
```

## Design
Follows SHAE Field Manual: warm paper #F2EEE5, ink #191712, oxide #B3341B, 0px radius, 1px borders.
