package main

import (
	"fmt"

	"github.com/CodeHunt7/go-blog-aggregator/internal/config"
)

func main() {
	// Чтаем конфиг при запуске
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	
	// Устанавливаем наше имя пользователя
	err = cfg.SetUser("Mikhail")
	if err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		return
	}

	// Читаем заново и выводим
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}
	fmt.Println(cfg)
}