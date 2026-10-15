package WorkManager

import (
	"BatcherHTTP/Configs"
	"BatcherHTTP/Database"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Batcher struct {
	inbox   []Configs.Task
	size    int
	timeout time.Duration
	sink    *pgx.Conn
	done    <-chan struct{}
}

func NewBatch(conn *pgx.Conn, size int, timeout time.Duration) *Batcher {
	return &Batcher{
		inbox:   make([]Configs.Task, 0, size*4),
		size:    size,
		timeout: timeout,
		sink:    conn,
		done:    make(<-chan struct{}),
	}
}

func TaskReadier(newBatch *Batcher, rt Configs.Task) {
	rt.Handle = true
	newBatch.inbox = append(newBatch.inbox, rt)
	fmt.Println(newBatch.inbox)
}

func BatchTracker(ticker *time.Ticker, conn *pgx.Conn, size int, timeout time.Duration, ch <-chan *Configs.Task, ctx context.Context) {
	batch := NewBatch(conn, size, timeout)

	for {
		select {
		case readyTask := <-ch:
			TaskReadier(batch, *readyTask)
			if len(batch.inbox) >= size {
				err := Database.InsertInto(ctx, conn, batch.inbox)
				if err != nil {
					fmt.Println("Батч не был отправлен")
					break
				}
				fmt.Println("Батч успешно отправлен, CAPACITY LIMIT")
				batch = NewBatch(conn, size, timeout)
			}
		case <-ticker.C:
			if len(batch.inbox) == 0 {
				fmt.Println("Очередь пустая")
				break
			}
			err := Database.InsertInto(ctx, conn, batch.inbox)
			batch = NewBatch(conn, size, timeout)
			if err != nil {
				fmt.Println("Батч не был отправлен")
				break
			}
			fmt.Println("Батч успешно отправлен, TIMEOUT")
		}
	}
}
