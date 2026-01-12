package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime time.Time

type PageData struct {
	Time          string
	Hostname      string
	OS            string
	Arch          string
	MemoryAlloc   string
	NumGoroutine  int
	NumCPU        int
	Uptime        string
	EstimatedCost string // <--- NEW: Real-time cost tracking
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get the Container ID (Hostname)
	hostname, _ := os.Hostname()

	// Cost Calculation: t2.micro is approx $0.0116/hour
	duration := time.Since(startTime)
	hours := duration.Hours()
	cost := hours * 0.0116

	data := PageData{
		Time:          time.Now().Format("15:04:05 Mon, 02 Jan"),
		Hostname:      hostname,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		MemoryAlloc:   fmt.Sprintf("%d MB", m.Alloc/1024/1024),
		NumGoroutine:  runtime.NumGoroutine(),
		NumCPU:        runtime.NumCPU(),
		Uptime:        duration.Round(time.Second).String(),
		EstimatedCost: fmt.Sprintf("$%.4f", cost), // Format to 4 decimal places
	}

	tmpl, err := template.ParseFiles("static/home.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func main() {
	startTime = time.Now()

	// We only need the home route now. The "website" is gone.
	http.HandleFunc("/", homeHandler)

	// Serve static assets (CSS/JS if we add them later)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Cloud Commander is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}