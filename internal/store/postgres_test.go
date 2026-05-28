package store

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockPostgresStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sql mock: %v", err)
	}
	return &PostgresStore{DB: db}, mock, func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close database: %v", err)
		}
	}
}

func TestPostgresPing(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectPing()

	if err := s.Ping(); err != nil {
		t.Fatalf("expected ping to succeed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.EnsureTable(); err != nil {
		t.Fatalf("expected ensure table to succeed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPostgresEnsureTableError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").WillReturnError(errors.New("create failed"))

	if err := s.EnsureTable(); err == nil {
		t.Fatal("expected ensure table error")
	}
}

func TestPostgresGetAll(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99).
		AddRow(2, "Gadget", 14.50)
	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("expected get all to succeed: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}

func TestPostgresGetAllEmpty(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "name", "price"})
	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("expected get all to succeed: %v", err)
	}
	if products == nil || len(products) != 0 {
		t.Fatalf("expected empty slice, got %#v", products)
	}
}

func TestPostgresGetAllQueryError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").WillReturnError(errors.New("query failed"))

	if _, err := s.GetAll(); err == nil {
		t.Fatal("expected query error")
	}
}

func TestPostgresGetAllScanError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow("bad-id", "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").WillReturnRows(rows)

	if _, err := s.GetAll(); err == nil {
		t.Fatal("expected scan error")
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Widget", 9.99)
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WithArgs(1).WillReturnRows(rows)

	product, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("expected get by id to succeed: %v", err)
	}
	if product.ID != 1 || product.Name != "Widget" {
		t.Fatalf("unexpected product: %+v", product)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WithArgs(999).WillReturnError(sql.ErrNoRows)

	_, err := s.GetByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresGetByIDError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id").WithArgs(1).WillReturnError(errors.New("select failed"))

	_, err := s.GetByID(1)
	if err == nil {
		t.Fatal("expected select error")
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO products").WithArgs("Widget", 9.99).WillReturnRows(rows)

	product, err := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	if product.ID != 1 {
		t.Fatalf("expected id 1, got %d", product.ID)
	}
}

func TestPostgresCreateError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectQuery("INSERT INTO products").WithArgs("Widget", 9.99).WillReturnError(errors.New("insert failed"))

	if _, err := s.Create(model.Product{Name: "Widget", Price: 9.99}); err == nil {
		t.Fatal("expected create error")
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("UPDATE products SET").WithArgs("Widget", 9.99, 1).WillReturnResult(sqlmock.NewResult(0, 1))

	product, err := s.Update(1, model.Product{Name: "Widget", Price: 9.99})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if product.ID != 1 {
		t.Fatalf("expected id 1, got %d", product.ID)
	}
}

func TestPostgresUpdateExecError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("UPDATE products SET").WithArgs("Widget", 9.99, 1).WillReturnError(errors.New("update failed"))

	if _, err := s.Update(1, model.Product{Name: "Widget", Price: 9.99}); err == nil {
		t.Fatal("expected update error")
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("UPDATE products SET").WithArgs("Missing", 1.00, 999).WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := s.Update(999, model.Product{Name: "Missing", Price: 1})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("DELETE FROM products WHERE id").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := s.Delete(1); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
}

func TestPostgresDeleteExecError(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("DELETE FROM products WHERE id").WithArgs(1).WillReturnError(errors.New("delete failed"))

	if err := s.Delete(1); err == nil {
		t.Fatal("expected delete error")
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()
	mock.ExpectExec("DELETE FROM products WHERE id").WithArgs(999).WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.Delete(999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
