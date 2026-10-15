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
	done    <-chan struct{}
}

func NewBatch(size int, timeout time.Duration) *Batcher {
	return &Batcher{
		inbox:   make([]Configs.Task, 0, size*4),
		size:    size,
		timeout: timeout,
		done:    make(<-chan struct{}),
	}
}

func TaskReadier(newBatch *Batcher, rt Configs.Task) {
	rt.Handle = true
	newBatch.inbox = append(newBatch.inbox, rt)
	//fmt.Println(newBatch.inbox)
}

func BatchTracker(ticker *time.Ticker, conn *pgx.Conn, size int, timeout time.Duration, ch <-chan *Configs.Task, ctx context.Context) {
	batch := NewBatch(size, timeout)

	for {
		select {
		case readyTask := <-ch:
			TaskReadier(batch, *readyTask)
			if len(batch.inbox) >= size {
				err := Database.InsertInto(ctx, conn, batch.inbox)
				if err != nil {
					fmt.Println("Батч не был отправлен", err)
					break
				}
				fmt.Println("Батч успешно отправлен, CAPACITY LIMIT")
				batch = NewBatch(size, timeout)
			}
		case <-ticker.C:
			ticker.Reset(timeout)
			if len(batch.inbox) == 0 {
				fmt.Println("Очередь пустая")
				break
			}
			err := Database.InsertInto(ctx, conn, batch.inbox)
			batch = NewBatch(size, timeout)
			if err != nil {
				fmt.Println("Батч не был отправлен", err)
				break
			}
			fmt.Println("Батч успешно отправлен, TIMEOUT")
		case <-ctx.Done():
			ticker.Stop()

			ctxNew, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if len(batch.inbox) == 0 {
				fmt.Println("Батч успешно отправлен, CONTEXT IS OVER")
				return
			}
			err := Database.InsertInto(ctxNew, conn, batch.inbox)
			if err != nil {
				fmt.Println("Батч не был отправлен", err)
				break
			}
			fmt.Println("Батч успешно отправлен, CONTEXT IS OVER")
			return
		}
	}
}
