package Database

import (
	"BatcherHTTP/Configs"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func InitDB(conn *pgx.Conn, ctx context.Context) error {
	createReq := `
		CREATE TABLE IF NOT EXISTS task_manager (
		id SERIAL PRIMARY KEY,
		info VARCHAR(256),
		received_at TIMESTAMP NOT NULL,
		handled BOOLEAN NOT NULL
		);
`
	_, err := conn.Exec(ctx, createReq)
	if err != nil {
		fmt.Println("ОШибка отправки SQL запроса")
		return err
	}
	return nil
}

func InsertInto(ctx context.Context, conn *pgx.Conn, rTasks []Configs.Task) error {
	//sqlRec := `INSERT INTO task_manager (info, received_at, handled) VALUES ($1, $2, $3);`

	amount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"task_manager"},
		[]string{"info", "received_at", "handled"},
		pgx.CopyFromSlice(len(rTasks), func(i int) ([]any, error) {
			return []any{
				rTasks[i].Info,
				rTasks[i].ReceivedAt,
				rTasks[i].Handle,
			}, nil
		}),
	)

	if err != nil {
		fmt.Println("Не удалось вставить батч:", err)
		return err
	}

	fmt.Printf("Успешно сброшено в БД: %d задач", amount)
	return nil
}
