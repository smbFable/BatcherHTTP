package main

import (
	"BatcherHTTP/Configs"
	"BatcherHTTP/Database"
	"BatcherHTTP/Database/WorkManager"
	"BatcherHTTP/Server"
	"BatcherHTTP/Server/Handlers"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const batchSize = 5
	const batchTimeout = 30 * time.Second
	chanTasks := make(chan *Configs.Task, 100)

	ticker := time.NewTicker(batchTimeout)

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
	fmt.Println("Подключение к Postgres успешно")

	err = Database.InitDB(conn, ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		_, err := Server.StartServer()
		if err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/tasks", Handlers.TaskHandler(chanTasks))

	go func() {
		WorkManager.BatchTracker(ticker, conn, batchSize, batchTimeout, chanTasks, ctx)
	}()

	<-ctx.Done()
}
