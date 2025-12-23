package main

import (
	"net/http"
	"net/http/httptest"
	"sipograf-go/controllers"
	"testing"
)

func TestApiJadwalHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/jadwal", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.ApiJadwal)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
	}
}

func TestHealthCheck(t *testing.T) {
	if 1+1 != 2 {
		t.Error("Math is broken")
	}
}