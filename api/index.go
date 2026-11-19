package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Item struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

var (
	items   = make(map[string]Item)
	nextID  int
	itemsMu sync.RWMutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", healthHandler)
	mux.HandleFunc("GET /api/data", apiDataHandler)
	mux.HandleFunc("POST /api/echo", echoHandler)
	mux.HandleFunc("GET /items", itemsListHandler)
	mux.HandleFunc("POST /items", itemsCreateHandler)
	mux.HandleFunc("GET /items/{id}", itemGetHandler)
	mux.HandleFunc("DELETE /items/{id}", itemDeleteHandler)
	
	log.Printf("SHAE Go HTTP Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "shae-go-http-server"})
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "SHAE Go HTTP Server", "version": "1.0.0"})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)
	json.NewEncoder(w).Encode(map[string]interface{}{"echo": body, "received_at": fmt.Sprintf("%d", os.Getpid())})
}

func itemsListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	itemsMu.RLock()
	list := make([]Item, 0, len(items))
	for _, v := range items {
		list = append(list, v)
	}
	itemsMu.RUnlock()
	json.NewEncoder(w).Encode(list)
}

func itemsCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	itemsMu.Lock()
	if item.ID == "" {
		nextID++
		item.ID = strconv.Itoa(nextID)
	}
	items[item.ID] = item
	itemsMu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": item.ID, "status": "created"})
}

func itemGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	itemsMu.RLock()
	item, ok := items[id]
	itemsMu.RUnlock()
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func itemDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	itemsMu.Lock()
	_, ok := items[id]
	if ok {
		delete(items, id)
	}
	itemsMu.Unlock()
	if !ok {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
