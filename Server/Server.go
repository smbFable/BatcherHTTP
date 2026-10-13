package Server

import (
	"BatcherHTTP/Configs"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func StartServer() (Configs.Server, error) {
	port := &Configs.Server{}

	data, err := os.ReadFile("Configs/config.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла")
		return Configs.Server{}, err
	}
	err = json.Unmarshal(data, port)
	if err != nil {
		fmt.Println("Ошибка чтения json")
		return Configs.Server{}, err
	}

	fmt.Println("Сервер запущен на порту", port.Port)

	err = http.ListenAndServe(port.Port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера")
		return Configs.Server{}, err
	}
	return Configs.Server{}, nil
}
