package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Item struct {
	ID    string `json:"id"`
	Data  string `json:"data"`
}

var items = make(map[string]Item)
var nextID int

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		list := make([]Item, 0, len(items))
		for _, v := range items {
			list = append(list, v)
		}
		json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if item.ID == "" {
			nextID++
			item.ID = strconv.Itoa(nextID)
		}
		items[item.ID] = item
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": item.ID, "status": "created"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		item, ok := items[id]
		if !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(item)
	case http.MethodDelete:
		if _, ok := items[id]; !ok {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		delete(items, id)
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"echo": true}`))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", healthHandler)
	mux.HandleFunc("GET /api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "SHAE Go HTTP Server"})
	})
	mux.HandleFunc("POST /api/echo", echoHandler)
	mux.HandleFunc("GET /items", itemsHandler)
	mux.HandleFunc("POST /items", itemsHandler)
	mux.HandleFunc("GET /items/{id}", itemHandler)
	mux.HandleFunc("DELETE /items/{id}", itemHandler)

	log.Printf("SHAE Go HTTP Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
