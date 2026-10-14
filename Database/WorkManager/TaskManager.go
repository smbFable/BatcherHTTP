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
		inbox:   make([]Configs.Task, size*4),
		size:    size,
		timeout: timeout,
		sink:    conn,
		done:    make(<-chan struct{}),
	}
}

func TaskReadier(since time.Time, newBatch *Batcher, rt Configs.Task) {
	time.Sleep(100 * time.Millisecond)

	rt.Handle = true
	rt.HandledAt = since

	newBatch.inbox = append(newBatch.inbox, rt)
}

func BatchTracker(ticker *time.Ticker, conn *pgx.Conn, size int, timeout time.Duration, ch <-chan *Configs.Task, since time.Time, ctx context.Context) {
	batch := NewBatch(conn, size, timeout)

	for {
		select {
		case readyTask := <-ch:
			TaskReadier(since, batch, *readyTask)
		case <-ticker.C:
			err := Database.InsertInto(ctx, conn, batch.inbox)
			batch = NewBatch(conn, size, timeout)
			if err != nil {
				fmt.Println("Батч не был отправлен")
				return
			}
			fmt.Println("Батч успешно отправлен")
		}
	}
}
