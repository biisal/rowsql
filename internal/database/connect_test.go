package database

// import (
// 	"context"
// 	"testing"
//
// 	"github.com/testcontainers/testcontainers-go/modules/postgres"
// )
//
// func TestDatabaseConnection(t *testing.T) {
// 	ctx := context.Background()
//
// 	pg, err := postgres.Run(ctx, "postgres:16-alpine", postgres.WithDatabase("test"),
// 		postgres.WithUsername("user"),
//
// 		postgres.WithPassword("password"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	connString, err := pg.ConnectionString(ctx, "sslmode=disable")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	t.Run("TestDatabaseConnection", func(t *testing.T) {
// 		_, err := CreateDBPool(ctx, connString)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})
// }
