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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	const batchSize = 5
	const batchTimeout = 30 * time.Second
	chanTasks := make(chan *Configs.Task, 100)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	wg := &sync.WaitGroup{}

	ticker := time.NewTicker(batchTimeout)

	pgxS, err := Database.ParserDB()
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(ctx, pgxS.Path)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Подключение к Postgres успешно")

	err = Database.InitDB(conn, ctx)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/tasks", Handlers.TaskHandler(chanTasks))

	go func() {
		_, err := Server.StartServer()
		if err != nil {
			panic(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		WorkManager.BatchTracker(ticker, conn, batchSize, batchTimeout, chanTasks, ctx)
	}()

	<-ctx.Done()

	wg.Wait()
	fmt.Println("Завершение работы")
}
