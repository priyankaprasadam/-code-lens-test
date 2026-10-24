package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	mu    = sync.Mutex{}
	tasks = []Task{{ID: 1, Title: "Build PR reviewer", Done: false}}
)

func getTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Unable to encode tasks", http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	mu.Lock()
	task.ID = len(tasks) + 1
	tasks = append(tasks, task)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("encode created task response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	//APIs
	mux.HandleFunc("GET /tasks", getTasks)
	mux.HandleFunc("POST /tasks", createTask)
	log.Println("Task API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
