package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sankalp-r/url-shortner/pkg/storage"
)

func TestShortenURL(t *testing.T) {
	store := storage.NewStore()
	handler := &Handler{store: store}

	router := mux.NewRouter()
	router.HandleFunc("/short", handler.ShortenURL).Methods("POST")
	router.HandleFunc("/{shortURL}", handler.RedirectURL).Methods("GET")

	reqBody, _ := json.Marshal(ShortenRequest{URL: "https://test.com"})
	req, err := http.NewRequest("POST", "/short", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res ShortenResponse
	err = json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	if res.ShortURL == "" {
		t.Errorf("handler returned an empty short URL")
	}

	req, err = http.NewRequest("GET", "/"+res.ShortURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	location := rr.Header().Get("Location")
	if location != "https://test.com" {
		t.Errorf("handler returned wrong location: got %v want %v", location, "https://example.com")
	}

}
