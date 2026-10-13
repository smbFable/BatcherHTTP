package main

import (
	"BatcherHTTP/Database"
	"BatcherHTTP/Server"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgxS, err := Database.ParserDB()
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(ctx, pgxS.Path)
	if err != nil {
		panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Подключение успешно")

	go func() {
		_, err := Server.StartServer()
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Контекст отменен, сервер остановлен")
}
