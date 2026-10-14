package Handlers

import (
	"BatcherHTTP/Configs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func TaskHandler(ch chan<- *Configs.Task, start time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			_, _ = fmt.Fprint(w, "Неверный метод запроса ", http.StatusBadRequest)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(r.Body)

		task := &Configs.Task{}

		byteData, err := io.ReadAll(r.Body)
		if err != nil {
			_, err = fmt.Fprint(w, "Ошибка body ", http.StatusNoContent)
			return
		}

		err = json.Unmarshal(byteData, task)
		task.ReceivedAt = time.Since(start)

		if err != nil {
			_, err = fmt.Fprint(w, "Ошибка чтения json ", http.StatusNoContent)
			return
		}
		_, err = fmt.Fprint(w, "Задание прочитано и отправлено на обработку ", http.StatusAccepted)

		ch <- task
	}
}
