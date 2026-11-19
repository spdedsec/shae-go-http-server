# Go HTTP Server

A serverless HTTP microservice in Go deployed on Vercel, featuring a built-in Apple Human Interface Guidelines (HIG) API Explorer and RESTful JSON endpoints.

## Features

- **Apple HIG API Explorer** — Cupertino-style interactive testing dashboard when requested in the browser
- **Serverless Architecture** — Vercel Go runtime entrypoint via `api/index.go` (`func Handler(w http.ResponseWriter, r *http.Request)`)
- **RESTful Endpoints** — Full CRUD support for `/items` with thread-safe in-memory storage (`sync.RWMutex`)
- **Proper Status Codes** — 200 OK, 201 Created, 400 Bad Request, 404 Not Found, 405 Method Not Allowed
- **CORS Enabled** — Configured for cross-origin browser and native client access
- **Light & Dark Mode** — Embedded web explorer with `color-scheme: light dark` and San Francisco typography

## Endpoints

| Method | Route | Description | Response |
|--------|-------|-------------|----------|
| `GET` | `/` | Apple HIG Web UI (Browser) / System Status (JSON) | HTML / 200 JSON |
| `GET` | `/health` | Health check endpoint | 200 JSON |
| `GET` | `/api/data` | Server metadata & version | 200 JSON |
| `POST` | `/api/echo` | JSON payload echo | 200 JSON |
| `GET` | `/items` | List all items | 200 JSON |
| `POST` | `/items` | Create a new item (`{"data":"..."}`) | 201 JSON |
| `GET` | `/items/{id}` | Retrieve item by ID | 200 / 404 JSON |
| `DELETE` | `/items/{id}` | Delete item by ID | 200 / 404 JSON |

## Local Development

```bash
go run api/index.go
```

Then visit `http://localhost:8080`.

## Deployment

Configured for Vercel with `"framework": "go"` and wildcard rewrites to `/api/index.go`.

## License

MIT © 2026 spdedsec
