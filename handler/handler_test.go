package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFaviconHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FaviconHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK && status != http.StatusNoContent {
		t.Errorf("FaviconHandler returned wrong status code: got %v", status)
	}
}
