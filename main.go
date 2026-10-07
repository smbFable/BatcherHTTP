package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pgxS := "postgres://postgres:A1234@localhost:5432/batcher"

	conn, err := pgx.Connect(ctx, pgxS)
	if err != nil {
		panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Подключение успешно")
}
