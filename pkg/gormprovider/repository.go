package gormprovider

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	NewQuery(ctx context.Context) *gorm.DB
	RunInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type RepositoryImp struct {
	db        *gorm.DB
	tableName string
}

func (r *RepositoryImp) NewQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table(r.tableName)
}

func (r *RepositoryImp) RunInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, contextTransactionKey, tx))
	})
}
