package db

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestTransaction_Begin(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tm := NewTransactionManager(db)
	ctx := context.Background()

	tx, err := tm.Begin(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tx == nil {
		t.Fatalf("expected transaction, got nil")
	}
}

func TestTransaction_Commit(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectCommit()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tm := NewTransactionManager(db)
	ctx := context.Background()

	_, err = tm.Begin(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = tm.Commit(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTransaction_Commit_NoActiveTransaction(t *testing.T) {
	mockDb, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tm := NewTransactionManager(db)
	ctx := context.Background()

	err = tm.Commit(ctx)
	if !errors.Is(err, ErrNoActiveTransactionToCommit) {
		t.Fatalf("expected %v, got %v", ErrNoActiveTransactionToCommit, err)
	}
}

func TestTransaction_Rollback(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectRollback()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tm := NewTransactionManager(db)
	ctx := context.Background()

	_, err = tm.Begin(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = tm.Rollback(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestTransaction_Rollback_NoActiveTransaction(t *testing.T) {
	mockDb, _, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	tm := NewTransactionManager(db)
	ctx := context.Background()

	err = tm.Rollback(ctx)
	if !errors.Is(err, ErrNoActiveTransactionToRollback) {
		t.Fatalf("expected ErrNoActiveTransactionToRollback, got %v", err)
	}
}
