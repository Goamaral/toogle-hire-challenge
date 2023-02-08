package repository_test

import (
	"os"
	"sync"
	"testing"
)

var databaseInitSql string
var databaseInitSqlMux sync.Once

func getDatabaseInitSql(t *testing.T) string {
	databaseInitSqlMux.Do(func() {
		databaseInitSqlBytes, err := os.ReadFile("../../database_init.sql")
		if err != nil {
			t.Fatal(err)
		}
		databaseInitSql = string(databaseInitSqlBytes)
	})
	return databaseInitSql
}
