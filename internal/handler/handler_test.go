package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupRouter() (*mux.Router, *Handler) {
	s := store.NewMemoryStore()
	h := NewHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, h
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductsEmpty(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	r, _ := setupRouter()

	// Create
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	// Get
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
	r, _ := setupRouter()

	createReq := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	createRR := httptest.NewRecorder()
	r.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createRR.Code)
	}

	updateReq := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Updated Widget","price":19.99}`))
	updateRR := httptest.NewRecorder()
	r.ServeHTTP(updateRR, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", updateRR.Code)
	}

	var updated model.Product
	if err := json.NewDecoder(updateRR.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}

	if updated.ID != 1 {
		t.Errorf("expected ID 1, got %d", updated.ID)
	}
	if updated.Name != "Updated Widget" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
	if updated.Price != 19.99 {
		t.Errorf("expected updated price 19.99, got %v", updated.Price)
	}
}

func TestDeleteProduct(t *testing.T) {
	r, _ := setupRouter()

	createReq := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	createRR := httptest.NewRecorder()
	r.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createRR.Code)
	}

	deleteReq := httptest.NewRequest("DELETE", "/products/1", nil)
	deleteRR := httptest.NewRecorder()
	r.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", deleteRR.Code)
	}

	getReq := httptest.NewRequest("GET", "/products/1", nil)
	getRR := httptest.NewRecorder()
	r.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getRR.Code)
	}
}

func TestCreateInvalidProduct(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
