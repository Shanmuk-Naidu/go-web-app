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
	
	// 3. Create the handler (UPDATED: using 'homeHandler' now)
	handlerFunc := http.HandlerFunc(homeHandler)

	// 4. Send the fake request
	handlerFunc.ServeHTTP(rr, req)

	// 5. Check if the status code is 200 (OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}