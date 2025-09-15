package main

import (
	"encoding/json" // for encoding/decoding JSON (request/response)
	"fmt"           // for console logs
	"net/http"      // for HTTP server + handlers
	"strconv"       // for converting strings to ints (task indexes in URL)
	"strings"       // for splitting URL paths like /lists/work/tasks/1
	"sync"          // for concurrency-safe map access (sync.Mutex)
)

// Task represents a single to-do item
type Task struct {
	Description string `json:"description"` // JSON tag lets API use {"description": "..."}
	Status      bool   `json:"status"`      // true = completed, false = incomplete
}

// Global shared state: map of task lists
// map[listName] -> []Task
var (
	taskLists = make(map[string][]Task) // initialized with make (empty map, ready to use)
	mu        sync.Mutex                // Mutex to avoid concurrent map writes
)

func main() {
	fmt.Println("Starting To-Do List HTTP server on http://localhost:8080")

	// Register HTTP handlers
	http.HandleFunc("/lists", handleLists)           // GET /lists -> list all lists
	http.HandleFunc("/lists/", handleListOperations) // operations on individual lists

	// Start HTTP server on port 8080
	// (blocks forever unless server crashes)
	http.ListenAndServe(":8080", nil)
}

// Handle "/lists" route
// GET /lists -> return all task lists
func handleLists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // Only GET allowed here
		mu.Lock()               // lock shared map
		defer mu.Unlock()       // ensure unlock even if function exits early
		writeJSON(w, taskLists) // send map as JSON response
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle "/lists/{listName}" and deeper paths like tasks
func handleListOperations(w http.ResponseWriter, r *http.Request) {
	// Trim the prefix and split path into parts
	// Example: "/lists/work/tasks/1" -> ["work", "tasks", "1"]
	path := strings.TrimPrefix(r.URL.Path, "/lists/")
	parts := strings.Split(path, "/")

	// Require at least a list name
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "List name required", http.StatusBadRequest)
		return
	}
	listName := parts[0]

	switch {
	// POST /lists/{listName} -> create new task list
	case len(parts) == 1 && r.Method == http.MethodPost:
		mu.Lock()
		defer mu.Unlock()
		if _, exists := taskLists[listName]; exists {
			http.Error(w, "Task list already exists", http.StatusConflict)
			return
		}
		taskLists[listName] = []Task{} // initialize empty slice
		writeJSON(w, map[string]string{"message": "Task list created", "list": listName})

	// GET /lists/{listName}/tasks -> return all tasks
	case len(parts) == 2 && parts[1] == "tasks" && r.Method == http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		tasks, exists := taskLists[listName]
		if !exists {
			http.Error(w, "Task list not found", http.StatusNotFound)
			return
		}
		writeJSON(w, tasks)

	// POST /lists/{listName}/tasks -> add new task (JSON body)
	case len(parts) == 2 && parts[1] == "tasks" && r.Method == http.MethodPost:
		var t Task
		// Decode JSON body into Task struct
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid task data", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		taskLists[listName] = append(taskLists[listName], t) // append new task
		writeJSON(w, map[string]string{"message": "Task added", "task": t.Description})

	// PUT /lists/{listName}/tasks/{index} -> toggle task status
	case len(parts) == 3 && parts[1] == "tasks" && r.Method == http.MethodPut:
		index, err := strconv.Atoi(parts[2]) // convert "1" -> 1
		if err != nil || index < 1 {
			http.Error(w, "Invalid task index", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		tasks, exists := taskLists[listName]
		if !exists {
			http.Error(w, "Task list not found", http.StatusNotFound)
			return
		}
		if index > len(tasks) {
			http.Error(w, "Task index out of range", http.StatusBadRequest)
			return
		}
		// Flip boolean (completed <-> incomplete)
		tasks[index-1].Status = !tasks[index-1].Status
		taskLists[listName] = tasks
		writeJSON(w, map[string]string{"message": "Task status toggled"})

	// DELETE /lists/{listName}/tasks/{index} -> delete task
	case len(parts) == 3 && parts[1] == "tasks" && r.Method == http.MethodDelete:
		index, err := strconv.Atoi(parts[2])
		if err != nil || index < 1 {
			http.Error(w, "Invalid task index", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		tasks, exists := taskLists[listName]
		if !exists {
			http.Error(w, "Task list not found", http.StatusNotFound)
			return
		}
		if index > len(tasks) {
			http.Error(w, "Task index out of range", http.StatusBadRequest)
			return
		}
		// Remove item at index: tasks[:i] + tasks[i+1:]
		taskLists[listName] = append(tasks[:index-1], tasks[index:]...)
		writeJSON(w, map[string]string{"message": "Task deleted"})

	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// Helper function: writes any Go value as JSON response
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data) // marshals data -> JSON and writes it
}
