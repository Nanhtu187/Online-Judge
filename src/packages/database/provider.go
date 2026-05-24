package database

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IProvider interface {
	Transact(ctx context.Context, fn func(ctx context.Context) error) error
	Readonly(ctx context.Context) context.Context
	GetDB() *gorm.DB
}

type Provider struct {
	db *gorm.DB
}

func NewProvider(db *gorm.DB) IProvider {
	return &Provider{db: db}
}

func (p *Provider) GetDB() *gorm.DB {
	return p.db
}

func (p *Provider) Transact(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	_, ok := getTxFromContext(ctx)
	if ok {
		return fn(ctx)
	}

	tx := p.db.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = tx.Rollback().Error
			zap.L().Error("panic recovered during transaction", zap.Any("panic", r))
			err = fmt.Errorf("panic during transaction: %v", r)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	ctx = context.WithValue(ctx, ctxTxKey, tx)

	if err = fn(ctx); err != nil {
		return err
	}

	return tx.Commit().Error
}

func (p *Provider) Readonly(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxReadonlyKey, p.db)
}

func getTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(ctxTxKey).(*gorm.DB)
	return tx, ok
}

// GetTx get Transaction from context
func GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := getTxFromContext(ctx)
	if !ok {
		err := errors.New("not found transaction in context")
		zap.L().Error(err.Error())
		return nil, err
	}
	return tx, nil
}

// GetReadonly get Readonly from context
func GetReadonly(ctx context.Context) (*gorm.DB, error) {
	tx, ok := getTxFromContext(ctx)
	if ok {
		return tx, nil
	}

	db, ok := ctx.Value(ctxReadonlyKey).(*gorm.DB)
	if ok {
		return db, nil
	}

	err := errors.New("not found readonly repository in context")
	zap.L().Error(err.Error())
	return nil, err
}

type ctxTxKeyType struct{}
type ctxReadonlyKeyType struct{}

var (
	ctxTxKey       = ctxTxKeyType{}
	ctxReadonlyKey = ctxReadonlyKeyType{}
)
