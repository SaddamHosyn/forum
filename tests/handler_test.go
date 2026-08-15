package tests

import (
	"forum-go/handler"
	"forum-go/render"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure working directory is project root
	if _, err := os.Stat("templates"); err != nil {
		_ = os.Chdir("..")
	}
	render.InitTemplates()
	os.Exit(m.Run())
}

func TestFaviconHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.FaviconHandler(rr, req)

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
	handler.ErrorHandler(rr, req, http.StatusNotFound)

	if rr.Code != http.StatusNotFound {
		t.Errorf("ErrorHandler returned status code: got %v, want %v", rr.Code, http.StatusNotFound)
	}
}
