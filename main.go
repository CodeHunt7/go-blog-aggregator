package main

import (
	"fmt"
	"os"

	"github.com/CodeHunt7/go-blog-aggregator/internal/config"
)

// Делаем структуру-рюкзак для всех данных (функции, бд) передаваеммых в команды
type state struct {
	cfg *config.Config
}

// Cтруктура для CLI-команды
type command struct {
	name string
	agrs []string
}

// Login handler function
func handlerLogin(s *state, cmd command) error {
	// Нет логина - ошибка
	if len(cmd.agrs) == 0 {
		return fmt.Errorf("Username is required")
	}

	// Устанавливаем наше имя пользователя
	err := s.cfg.SetUser(cmd.agrs[0])
	if err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		return err
	}

	fmt.Printf("User %s has been set", s.cfg.CurrentUserName)
	return nil
}

// Словарь всех CLI-команд
type commands struct {
	handlerMap map[string]func(*state, command) error
}

// Метод запускающий команду с параметрами (если есть)
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlerMap[cmd.name]
	if !exists {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

// Метод регистрирующий хэндлер-функцию для команды
func (c *commands) register(name string, f func(*state, command) error) {
	_, exists := c.handlerMap[name]
	if exists {
		fmt.Printf("Command %s already exists\n", name)
		return
	}
	c.handlerMap[name] = f
} 

func main() {
	// Чтаем конфиг при запуске
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Инициализируем состояние и команды
	var s state	
	s.cfg = &cfg
	var c commands
	c.handlerMap = make(map[string]func(*state, command) error)

	// Регистрируем команду login
	c.register("login", handlerLogin)

	// Берем аргументы из ввода пользователя
	if len(os.Args) < 2 {
        fmt.Println("Not enough arguments")
        os.Exit(1) // Выход с ошибкой
    }
	cmd := command{
        name: os.Args[1],
        agrs: os.Args[2:],
    }

	// Запускаем команду
	err = c.run(&s, cmd)
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1) // Выход с ошибкой
	}

}