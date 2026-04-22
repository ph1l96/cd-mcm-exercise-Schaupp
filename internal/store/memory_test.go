package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	// create a product and verify GetByID returns it
	testProduct := s.Create(model.Product{Name: "Beer", Price: 5.20})

	if testProduct.ID == 0 {
		t.Fatal(("Expected testProduct to have an ID >0"))
	}

	got, err := s.GetByID(testProduct.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if got.ID != testProduct.ID {
		t.Errorf("Expected id %d, got %d", testProduct.ID, got.ID)
	}

	if got.Name != testProduct.Name {
		t.Errorf("Expected name %q, got %q", testProduct.Name, got.Name)
	}

	if got.Price != testProduct.Price {
		t.Errorf("Expected price %v, got %v", testProduct.Price, got.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	testProduct := s.Create(model.Product{Name: "Beer", Price: 5.20})

	err := s.Delete(testProduct.ID)
	if err != nil {
		t.Fatalf("expected no error from Delete, got %v", err)
	}

	_, err = s.GetByID(testProduct.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after deleting product, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound for non-existent product, got %v", err)
	}
}

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()
	testProduct := s.Create(model.Product{Name: "Beer", Price: 5.20})

	updatedProduct, err := s.Update(testProduct.ID, model.Product{Name: "Beer", Price: 4.99})
	if err != nil {
		t.Fatalf("expected no error from Update, got %v", err)
	}

	got, err := s.GetByID(testProduct.ID)
	if err != nil {
		t.Fatalf("expected no error from GetByID, got %v", err)
	}

	if got.ID != testProduct.ID {
		t.Errorf("expected ID %d, got %d", testProduct.ID, got.ID)
	}

	if got.Name != updatedProduct.Name {
		t.Errorf("expected name %q, got %q", updatedProduct.Name, got.Name)
	}

	if got.Price != updatedProduct.Price {
		t.Errorf("expected price %v, got %v", updatedProduct.Price, got.Price)
	}
}

func TestProductValidate(t *testing.T) {
	tests := []struct {
		name     string
		product  model.Product
		expected bool
	}{
		{
			name:     "empty name",
			product:  model.Product{Name: "", Price: 10.0},
			expected: false,
		},
		{
			name:     "negative price",
			product:  model.Product{Name: "Widget", Price: -5.0},
			expected: false,
		},
		{
			name:     "valid product",
			product:  model.Product{Name: "Widget", Price: 9.99},
			expected: true,
		},
	}

	for _, tt := range tests {
		if tt.product.Validate() != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, tt.product.Validate())
		}
	}
}
