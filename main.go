package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

// Define a struct to hold the data for the Dashboard
type PageData struct {
	Time         string
	OS           string
	Arch         string
	MemoryAlloc  string
	NumGoroutine int
}

// Handler for the Dashboard (Home Page)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// If the URL path is NOT exactly "/", serve a 404 (or handle specific files)
	// This prevents "/random" from loading the dashboard
	if r.URL.Path != "/" && r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	// 1. Get Memory Stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 2. Prepare the data
	data := PageData{
		Time:         time.Now().Format("2006-01-02 15:04:05"),
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		MemoryAlloc:  fmt.Sprintf("%d MB", m.Alloc/1024/1024),
		NumGoroutine: runtime.NumGoroutine(),
	}

	// 3. Parse and execute the template
	tmpl, err := template.ParseFiles("static/home.html")
	if err != nil {
		http.Error(w, "Could not load home template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Helper function to serve simple static pages
func staticHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/"+filename)
	}
}

func main() {
	// Route 1: The Dashboard (Home)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)

	// Route 2: The Other Pages (Restoring your links)
	// Assuming your files are named about.html, contact.html, courses.html
	http.HandleFunc("/about", staticHandler("about.html"))
	http.HandleFunc("/contact", staticHandler("contact.html"))
	http.HandleFunc("/courses", staticHandler("courses.html"))

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}