package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// 1. Create a fake request to the home page
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a ResponseRecorder (a fake browser)
	rr := httptest.NewRecorder()
	
	// 3. Create the handler (UPDATED: using 'handler' instead of 'homePage')
	handlerFunc := http.HandlerFunc(handler)

	// 4. Send the fake request
	handlerFunc.ServeHTTP(rr, req)

	// 5. Check if the status code is 200 (OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Optional: Check if the response contains the dashboard title
	// This confirms the template loaded correctly
	// expected := "DevOps System Monitor"
	// if !strings.Contains(rr.Body.String(), expected) {
	// 	t.Errorf("handler returned unexpected body: does not contain %v", expected)
	// }
}