package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

type fakePostgresStore struct {
	products []model.Product
	product  model.Product
	err      error
}

func (s *fakePostgresStore) Ping() error {
	return s.err
}

func (s *fakePostgresStore) GetAll() ([]model.Product, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.products, nil
}

func (s *fakePostgresStore) GetByID(id int) (model.Product, error) {
	if s.err != nil {
		return model.Product{}, s.err
	}
	return s.product, nil
}

func (s *fakePostgresStore) Create(p model.Product) (model.Product, error) {
	if s.err != nil {
		return model.Product{}, s.err
	}
	p.ID = 1
	return p, nil
}

func (s *fakePostgresStore) Update(id int, p model.Product) (model.Product, error) {
	if s.err != nil {
		return model.Product{}, s.err
	}
	p.ID = id
	return p, nil
}

func (s *fakePostgresStore) Delete(id int) error {
	return s.err
}

func setupPostgresRouter(s *fakePostgresStore) *mux.Router {
	h := NewPostgresHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestPostgresHealth(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHealthUnavailable(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: errors.New("down")})

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}
}

func TestPostgresGetProducts(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{
		products: []model.Product{{ID: 1, Name: "Widget", Price: 9.99}},
	})

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductsError(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: errors.New("query failed")})

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresGetProduct(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{product: model.Product{ID: 1, Name: "Widget", Price: 9.99}})

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: store.ErrNotFound})

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresCreateProduct(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidPayload(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductValidationError(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":1}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductStoreError(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: errors.New("insert failed")})

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresUpdateProduct(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Updated","price":11.50}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductInvalidPayload(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: store.ErrNotFound})

	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(`{"name":"Missing","price":1}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresDeleteProduct(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{})

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresDeleteProductNotFound(t *testing.T) {
	r := setupPostgresRouter(&fakePostgresStore{err: store.ErrNotFound})

	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
