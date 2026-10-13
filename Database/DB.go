package Database

import (
	"BatcherHTTP/Configs"
	"encoding/json"
	"fmt"
	"os"
)

func ParserDB() (Configs.DB, error) {
	data := &Configs.DB{}

	dataByte, err := os.ReadFile("Configs/config.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла")
		return Configs.DB{}, err
	}

	err = json.Unmarshal(dataByte, data)
	if err != nil {
		fmt.Println("Ошибка чтения байтов")
		return Configs.DB{}, err
	}
	return *data, nil
}
