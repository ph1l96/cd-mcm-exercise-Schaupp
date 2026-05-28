package store

import (
	"errors"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected product, got error: %v", err)
	}
	if got.ID != 1 || got.Name != "Widget" || got.Price != 9.99 {
		t.Fatalf("unexpected product: %+v", got)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetAllWithProducts(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "Widget", Price: 9.99})
	s.Create(model.Product{Name: "Gadget", Price: 14.50})

	products := s.GetAll()
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}

func TestUpdateExistingProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})

	updated, err := s.Update(created.ID, model.Product{Name: "Updated", Price: 12.34})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updated.ID != created.ID || updated.Name != "Updated" || updated.Price != 12.34 {
		t.Fatalf("unexpected updated product: %+v", updated)
	}
}

func TestUpdateNonExistentProduct(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, model.Product{Name: "Missing", Price: 1})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExistingProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if _, err := s.GetByID(created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted product to be missing, got %v", err)
	}
}
