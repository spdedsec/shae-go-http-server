package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

var (
	itemsMu sync.RWMutex
	items   = map[string]Item{
		"1": {ID: "1", Data: "Apple HIG System Integration"},
		"2": {ID: "2", Data: "Production Vercel Go Serverless Handler"},
	}
	nextID = 2
)

// Embedded Apple HIG Web UI for browser requests
const appleHIGHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="color-scheme" content="light dark">
  <meta name="description" content="Go HTTP Server - Apple HIG API Explorer">
  <title>Go HTTP Server &mdash; API Studio</title>
  <style>
    :root {
      color-scheme: light dark;
      --system-blue: #007AFF;
      --system-blue-hover: #0062CC;
      --system-green: #34C759;
      --system-red: #FF3B30;
      --system-orange: #FF9500;
      --system-bg: #FFFFFF;
      --system-secondary-bg: #F2F2F7;
      --system-card-bg: rgba(255, 255, 255, 0.88);
      --system-label: #000000;
      --system-secondary-label: rgba(60, 60, 67, 0.60);
      --system-border: rgba(0, 0, 0, 0.10);
      --system-fill: rgba(120, 120, 128, 0.16);
      --system-chrome: rgba(255, 255, 255, 0.82);
      --font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "SF Pro Display", "SF Pro", system-ui, sans-serif;
      --font-mono: -apple-system-ui-monospaced, "SF Mono", Menlo, Monaco, Consolas, monospace;
      --radius-sm: 8px;
      --radius-md: 12px;
      --radius-lg: 16px;
      --radius-pill: 9999px;
      --shadow-2: 0 4px 6px -1px rgba(0, 0, 0, 0.07), 0 2px 4px -1px rgba(0, 0, 0, 0.04);
      --spring: cubic-bezier(0.4, 0, 0.2, 1);
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --system-blue: #0A84FF;
        --system-blue-hover: #409CFF;
        --system-green: #30D158;
        --system-red: #FF453A;
        --system-orange: #FF9F0A;
        --system-bg: #000000;
        --system-secondary-bg: #1C1C1E;
        --system-card-bg: rgba(28, 28, 30, 0.88);
        --system-label: #FFFFFF;
        --system-secondary-label: rgba(235, 235, 245, 0.60);
        --system-border: rgba(255, 255, 255, 0.14);
        --system-fill: rgba(120, 120, 128, 0.36);
        --system-chrome: rgba(28, 28, 30, 0.82);
      }
    }
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: var(--font-family);
      background-color: var(--system-secondary-bg);
      color: var(--system-label);
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      -webkit-font-smoothing: antialiased;
    }
    .header {
      height: 56px;
      background: var(--system-chrome);
      backdrop-filter: blur(20px);
      -webkit-backdrop-filter: blur(20px);
      border-bottom: 1px solid var(--system-border);
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 24px;
      position: sticky;
      top: 0;
      z-index: 10;
    }
    .brand { display: flex; align-items: center; gap: 8px; }
    .badge {
      font-size: 11px;
      font-weight: 700;
      padding: 2px 8px;
      background: var(--system-blue);
      color: #fff;
      border-radius: var(--radius-pill);
    }
    .title { font-size: 17px; font-weight: 600; }
    .status-badge {
      font-size: 12px;
      font-weight: 600;
      color: var(--system-green);
      display: flex;
      align-items: center;
      gap: 6px;
    }
    .dot { width: 8px; height: 8px; border-radius: 50%; background: var(--system-green); }
    .container {
      max-width: 1000px;
      margin: 0 auto;
      width: 100%;
      padding: 32px 24px;
      display: flex;
      flex-direction: column;
      gap: 24px;
      flex: 1;
    }
    .card {
      background: var(--system-card-bg);
      backdrop-filter: blur(20px);
      -webkit-backdrop-filter: blur(20px);
      border-radius: var(--radius-lg);
      border: 1px solid var(--system-border);
      box-shadow: var(--shadow-2);
      padding: 24px;
    }
    .card-title {
      font-size: 14px;
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.04em;
      color: var(--system-secondary-label);
      margin-bottom: 16px;
    }
    .endpoint-grid { display: flex; flex-direction: column; gap: 12px; }
    .endpoint-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 12px 16px;
      background: var(--system-fill);
      border-radius: var(--radius-md);
      gap: 16px;
    }
    .method-tag {
      font-family: var(--font-mono);
      font-size: 12px;
      font-weight: 700;
      padding: 4px 8px;
      border-radius: var(--radius-sm);
      min-width: 54px;
      text-align: center;
    }
    .get { background: rgba(52, 199, 89, 0.2); color: var(--system-green); }
    .post { background: rgba(0, 122, 255, 0.2); color: var(--system-blue); }
    .delete { background: rgba(255, 59, 48, 0.2); color: var(--system-red); }
    .path { font-family: var(--font-mono); font-size: 14px; flex: 1; }
    .desc { font-size: 13px; color: var(--system-secondary-label); }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 44px;
      min-width: 44px;
      padding: 0 16px;
      font-family: var(--font-family);
      font-size: 13px;
      font-weight: 500;
      border-radius: var(--radius-pill);
      border: none;
      cursor: pointer;
      background: var(--system-blue);
      color: #fff;
      transition: all 150ms var(--spring);
    }
    .btn:hover { background: var(--system-blue-hover); transform: scale(0.98); }
    .btn:focus-visible { outline: 2px solid var(--system-blue); outline-offset: 2px; }
    .response-area {
      margin-top: 16px;
      background: var(--system-bg);
      border: 1px solid var(--system-border);
      border-radius: var(--radius-md);
      padding: 16px;
      font-family: var(--font-mono);
      font-size: 13px;
      max-height: 240px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
    .footer {
      height: 48px;
      padding: 0 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 12px;
      color: var(--system-secondary-label);
      border-top: 1px solid var(--system-border);
    }
    @media (prefers-reduced-motion: reduce) {
      *, *::before, *::after { animation-duration: 0.01ms !important; transition-duration: 0.01ms !important; }
    }
  </style>
</head>
<body>
  <header class="header">
    <div class="brand">
      <span class="badge">GO</span>
      <h1 class="title">HTTP Server</h1>
    </div>
    <div class="status-badge">
      <span class="dot"></span>
      <span>200 OK &middot; Vercel Serverless</span>
    </div>
  </header>

  <main class="container">
    <section class="card">
      <h2 class="card-title">API Endpoints (Apple HIG Explorer)</h2>
      <div class="endpoint-grid">
        <div class="endpoint-row">
          <span class="method-tag get">GET</span>
          <span class="path">/api/data</span>
          <span class="desc">Server metadata &amp; version</span>
          <button class="btn" onclick="runReq('GET', '/api/data')">Test</button>
        </div>
        <div class="endpoint-row">
          <span class="method-tag get">GET</span>
          <span class="path">/items</span>
          <span class="desc">List all items</span>
          <button class="btn" onclick="runReq('GET', '/items')">Test</button>
        </div>
        <div class="endpoint-row">
          <span class="method-tag post">POST</span>
          <span class="path">/items</span>
          <span class="desc">Create new item</span>
          <button class="btn" onclick="runReq('POST', '/items', {data: 'New item from Apple HIG Studio'})">Create</button>
        </div>
        <div class="endpoint-row">
          <span class="method-tag get">GET</span>
          <span class="path">/health</span>
          <span class="desc">System healthcheck</span>
          <button class="btn" onclick="runReq('GET', '/health')">Test</button>
        </div>
      </div>

      <div class="response-area" id="respOutput">// Click "Test" on any endpoint above to execute live request</div>
    </section>
  </main>

  <footer class="footer">
    <div>Apple Human Interface Guidelines &middot; Go Serverless Runtime</div>
    <div><a href="https://github.com/spdedsec/shae-go-http-server" target="_blank" rel="noopener" style="color:var(--system-blue); text-decoration:none;">GitHub</a></div>
  </footer>

  <script>
    async function runReq(method, path, body) {
      const out = document.getElementById('respOutput');
      out.textContent = 'Executing ' + method + ' ' + path + '...';
      try {
        const opts = { method, headers: { 'Accept': 'application/json' } };
        if (body) {
          opts.headers['Content-Type'] = 'application/json';
          opts.body = JSON.stringify(body);
        }
        const res = await fetch(path, opts);
        const data = await res.json();
        out.textContent = 'Status: ' + res.status + ' ' + res.statusText + '\n\n' + JSON.stringify(data, null, 2);
      } catch (err) {
        out.textContent = 'Error: ' + err.message;
      }
    }
  </script>
</body>
</html>`

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
}

// Handler is the Vercel serverless function entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path
	// Strip prefix if any
	path = strings.TrimPrefix(path, "/api/index")
	if path == "" {
		path = "/"
	}

	// Root path handler
	if path == "/" {
		accept := r.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(appleHIGHTML))
			return
		}
		// JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"service": "shae-go-http-server",
			"version": "1.0.0",
			"design":  "Apple HIG",
		})
		return
	}

	// Health route
	if path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "shae-go-http-server",
		})
		return
	}

	// API Data route
	if path == "/api/data" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Go HTTP Server &mdash; Apple HIG Studio",
			"version": "1.0.0",
			"runtime": "Vercel Go Serverless",
		})
		return
	}

	// Echo route
	if path == "/api/echo" {
		w.Header().Set("Content-Type", "application/json")
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"echo":        body,
			"received_at": fmt.Sprintf("%d", os.Getpid()),
		})
		return
	}

	// Items collection
	if path == "/items" {
		switch r.Method {
		case http.MethodGet:
			itemsMu.RLock()
			list := make([]Item, 0, len(items))
			for _, v := range items {
				list = append(list, v)
			}
			itemsMu.RUnlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(list)
			return

		case http.MethodPost:
			var item Item
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid json payload"})
				return
			}
			itemsMu.Lock()
			if item.ID == "" {
				nextID++
				item.ID = strconv.Itoa(nextID)
			}
			items[item.ID] = item
			itemsMu.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": item.ID, "status": "created"})
			return

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}
	}

	// Item by ID: /items/{id}
	if strings.HasPrefix(path, "/items/") {
		id := strings.TrimPrefix(path, "/items/")
		switch r.Method {
		case http.MethodGet:
			itemsMu.RLock()
			item, ok := items[id]
			itemsMu.RUnlock()
			w.Header().Set("Content-Type", "application/json")
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "item not found"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
			return

		case http.MethodDelete:
			itemsMu.Lock()
			_, ok := items[id]
			if ok {
				delete(items, id)
			}
			itemsMu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "item not found"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "deleted", "id": id})
			return

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}
	}

	// Fallback 404
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "route not found", "path": path})
}

// main entrypoint for local execution
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Go HTTP Server listening on http://0.0.0.0:%s\n", port)
	http.HandleFunc("/", Handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
