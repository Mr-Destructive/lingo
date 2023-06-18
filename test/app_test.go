package app_test

import (
	"lingo/lingo/app"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLinksHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/links", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.LinksHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `<a href="/add/link">Add Link</a>`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAddLinkHandler(t *testing.T) {
	// Test adding a new link
}

func TestEditLinkHandler(t *testing.T) {
	// Test editing an existing link
}

func TestDeleteLinkHandler(t *testing.T) {
	// Test deleting a link
}
