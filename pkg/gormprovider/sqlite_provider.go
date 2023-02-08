package gormprovider

import (
	"challenge/pkg/env"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

const contextTransactionKey ContextString = "tx"

type ContextString string

type SQLiteProvider struct {
	*gorm.DB
}

func NewSQLiteProvider() (*SQLiteProvider, error) {
	dbPath := env.GetOrDefault("SQLITE_DBPATH", "challenge.sqlite")
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbPath)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &SQLiteProvider{DB: db}, nil
}

func NewTestSQLiteProvider(t *testing.T, databaseInitSql string) *SQLiteProvider {
	dbFilename := fmt.Sprintf("test_%s.sqlite", ulid.Make().String())
	dbPath := strings.ToLower(fmt.Sprintf("%s?_pragma=foreign_keys(1)", dbFilename))
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(string(databaseInitSql)).Error
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(dbFilename)
	})

	return &SQLiteProvider{db}
}

func (p *SQLiteProvider) NewRepository(tableName string) Repository {
	return &RepositoryImp{db: p.DB, tableName: tableName}
}

func (p *SQLiteProvider) RunInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, contextTransactionKey, tx))
	})
}
