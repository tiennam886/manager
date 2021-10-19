package manager

import (
	"context"
	"time"

	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/go-sql-driver/mysql"
)

var (
	uri      = "mongodb://localhost:27017"
	database = "local"
	mySqlDB  *sql.DB
)

func connectCol(uri string, database string, col string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx := initCtx()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(database).Collection(col), nil
}

func initCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func connectMySql() (*sql.DB, error) {
	mySqlDB, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/app")
	return mySqlDB, err
}
