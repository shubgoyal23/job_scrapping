package helpers

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDataBase() {
	connStr := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), connStr)
	Conn = conn
	if err != nil {
		panic(err)
	}
}

func Insert(val map[string]bool) {
	batch := &pgx.Batch{}
	for i := range val {
		batch.Queue("INSERT INTO links (link) VALUES ($1)", i)
	}
	logs := Conn.SendBatch(context.Background(), batch)
	println(logs.Close().Error())
}
