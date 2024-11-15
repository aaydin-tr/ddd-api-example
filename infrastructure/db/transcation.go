package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type TransactionManager interface {
	Begin(ctx context.Context) (*gorm.DB, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

var (
	ErrNoActiveTransactionToCommit   = errors.New("no active transaction to commit")
	ErrNoActiveTransactionToRollback = errors.New("no active transaction to rollback")
)

type Transaction struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &Transaction{db: db}
}

func (t *Transaction) Begin(ctx context.Context) (*gorm.DB, error) {
	if t.tx != nil {
		return t.tx, nil
	}

	t.tx = t.db.Begin()
	if t.tx.Error != nil {
		return nil, t.tx.Error
	}

	return t.tx, nil
}

func (t *Transaction) Commit(ctx context.Context) error {
	if t.tx == nil {
		return ErrNoActiveTransactionToCommit
	}

	err := t.tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) Rollback(ctx context.Context) error {
	if t.tx == nil {
		return ErrNoActiveTransactionToRollback
	}

	err := t.tx.Rollback().Error
	if err != nil {
		return err
	}

	return nil
}
