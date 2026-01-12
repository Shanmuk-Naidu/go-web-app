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
	EstimatedCost string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	hostname, _ := os.Hostname()

	// Cost Calculation: t2.micro is approx $0.0116/hour
	duration := time.Since(startTime)
	hours := duration.Hours()
	cost := hours * 0.0116

	// FIX: Define the location BEFORE the struct
	loc := time.FixedZone("IST", 5.5*60*60) // +5:30 offset

	data := PageData{
		// FIX: Use .In(loc) to apply the timezone
		Time:          time.Now().In(loc).Format("15:04:05 Mon, 02 Jan"),
		Hostname:      hostname,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		MemoryAlloc:   fmt.Sprintf("%d MB", m.Alloc/1024/1024),
		NumGoroutine:  runtime.NumGoroutine(),
		NumCPU:        runtime.NumCPU(),
		Uptime:        duration.Round(time.Second).String(),
		EstimatedCost: fmt.Sprintf("$%.4f", cost),
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

	http.HandleFunc("/", homeHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Cloud Commander is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}