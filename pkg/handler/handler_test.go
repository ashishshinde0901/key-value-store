package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerlacdt/db-key-value-store/pkg/db"
	"github.com/mattetti/filebuffer"
)

func setup(t *testing.T) http.Handler {
	t.Parallel()
	h, err := New(db.New(filebuffer.New(nil)))
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}
	return h
}

func TestSingleHttpDelete(t *testing.T) {
	r := setup(t)
	srv := httptest.NewServer(r)

	// act
	value := []byte("bar")
	resp, err := http.Post(fmt.Sprintf("%s/db/foo", srv.URL), "application/octet-stream", bytes.NewReader(value))
	if err != nil {
		t.Fatalf("error http SET %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("statusCode expected %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/db/foo", srv.URL), nil)
	if err != nil {
		t.Fatalf("error creating DELETE request %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error DELETE %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("statusCode expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
	resp, err = http.Get(fmt.Sprintf("%s/db/foo", srv.URL))
	if err != nil {
		t.Fatalf("error http GET %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("statusCode expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestSingleHttpSetAndGet(t *testing.T) {
	r := setup(t)
	srv := httptest.NewServer(r)

	// act
	value := []byte("bar")
	resp, err := http.Post(fmt.Sprintf("%s/db/foo", srv.URL), "application/octet-stream", bytes.NewReader(value))
	if err != nil {
		t.Fatalf("error http SET %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("statusCode expected %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	resp, err = http.Get(fmt.Sprintf("%s/db/foo", srv.URL))
	if err != nil {
		t.Fatalf("error http GET %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("statusCode expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if string(body) != string(value) {
		t.Fatalf("body expected %s, got %s", value, body)
	}
}
