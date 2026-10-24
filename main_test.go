package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTasks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	getTasks(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got []Task
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one sample task")
	}
}

func TestCreateTask(t *testing.T) {
	mu.Lock()
	oldTasks := tasks
	tasks = []Task{}
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		tasks = oldTasks
		mu.Unlock()
	})
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Test task"}`))
	rec := httptest.NewRecorder()
	createTask(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	var got Task
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Title != "Test task" || got.ID != 1 {
		t.Fatalf("unexpected task: %#v", got)
	}
}

func TestCreateTaskRejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{"))
	rec := httptest.NewRecorder()
	createTask(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
