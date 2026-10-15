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
	sqlRec := `INSERT INTO task_manager (info, received_at, handled) VALUES ($1, $2, $3);`
	for _, rTask := range rTasks {
		info := rTask.Info
		receivedAt := rTask.ReceivedAt
		handled := rTask.Handle

		_, err := conn.Exec(ctx, sqlRec, info, receivedAt, handled)
		if err != nil {
			fmt.Println("Не удалось вставить строки", err)
			return err
		}
	}
	return nil
}
