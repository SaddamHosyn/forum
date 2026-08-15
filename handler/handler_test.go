package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure working directory is project root so relative paths (templates, assets) resolve correctly
	if _, err := os.Stat("../templates"); err == nil {
		_ = os.Chdir("..")
	}
	os.Exit(m.Run())
}

func TestFaviconHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	FaviconHandler(rr, req)

	// Status can be OK (200) or NotFound (404) depending on asset presence
	if rr.Code != http.StatusOK && rr.Code != http.StatusNotFound {
		t.Errorf("FaviconHandler returned status code: got %v", rr.Code)
	}
}

func TestErrorHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ErrorHandler(rr, req, http.StatusNotFound)

	if rr.Code != http.StatusNotFound {
		t.Errorf("ErrorHandler returned status code: got %v, want %v", rr.Code, http.StatusNotFound)
	}
}
