package main

import (
    "fmt"
    "html/template"
    "net/http"
    "runtime"
    "time"
)

// Define a struct to hold the data we send to the HTML page
type PageData struct {
    Time         string
    OS           string
    Arch         string
    MemoryAlloc  string
    NumGoroutine int
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Get Memory Stats
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    // 2. Prepare the data
    data := PageData{
        Time:         time.Now().Format("2006-01-02 15:04:05"),
        OS:           runtime.GOOS,
        Arch:         runtime.GOARCH,
        // Convert bytes to Megabytes (MB)
        MemoryAlloc:  fmt.Sprintf("%d MB", m.Alloc/1024/1024),
        NumGoroutine: runtime.NumGoroutine(),
    }

    // 3. Parse the template (UPDATED: pointing to home.html)
    tmpl, err := template.ParseFiles("static/home.html")
    if err != nil {
        http.Error(w, "Could not load template", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}