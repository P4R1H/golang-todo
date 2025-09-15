# Learning Backend Development with Go
Notes while building a to-do list backend in Go.

---

## 1. Routing in Go’s stdlib
- `http.ServeMux` uses **prefix-based routing**.
- No wildcards or `{param}` support by default.
- To extract data from the path, you do string operations manually.

Example:

```go
path := strings.TrimPrefix(r.URL.Path, "/lists/")
parts := strings.Split(path, "/")

listName := parts[0] // e.g. "work"
taskNum  := parts[1] // e.g. "3"
````

Request:

```
GET /lists/work/tasks/3
```

Result:

```
parts = ["work", "tasks", "3"]
```

If you want `/lists/{listname}/tasks/{id}` style, use a third-party router like `chi` or `gorilla/mux`.

---

## 2. Writing responses

* Handlers **must write to the ResponseWriter (`w`)**.
* You cannot `return` a value directly.

Example:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World")
}
```

---

## 3. JSON responses

* Default output is plain text.
* To return JSON, set headers and encode:

```go
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

Usage:

```go
writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
```

---

## 4. Checking existence in maps

* Idiomatic Go pattern:

```go
if _, exists := taskLists[listName]; exists {
    fmt.Println("Task list already exists")
} else {
    taskLists[listName] = []Task{}
}
```

---

## 5. Common `http` functions

* `http.HandleFunc(pattern, handler)` → register routes
* `http.ListenAndServe(":8080", nil)` → start server
* `r.URL.Path` → request path
* `r.Method` → HTTP method (`GET`, `POST`, etc.)
* `r.Body` → request body (use `json.NewDecoder(r.Body)` to parse JSON)

---

## 6. Typical project layout (for later)

* `main.go` → entry point
* `handlers.go` → request handlers
* `models.go` → data structures
* `utils.go` → helper functions (e.g., `writeJSON`)

---