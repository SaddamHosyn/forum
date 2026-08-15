package tests

import (
	"bytes"
	"forum-go/handler"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandlerMethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest("GET", "/register", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.RegisterHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("RegisterHandler GET status code: got %v, want %v", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestLoginHandlerInvalidUser(t *testing.T) {
	_ = setupTestDB(t)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("username", "nonexistentuser")
	_ = w.WriteField("password", "SomePass123!")
	w.Close()

	req, err := http.NewRequest("POST", "/login", &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	rr := httptest.NewRecorder()
	handler.LoginHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("LoginHandler invalid user status: got %v, want %v", rr.Code, http.StatusUnauthorized)
	}
}
