package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/CodeHunt7/go-blog-aggregator/internal/config"
	"github.com/CodeHunt7/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Делаем структуру-рюкзак для всех данных (функции, бд) передаваеммых в команды
type state struct {
	cfg *config.Config
	db  *database.Queries
}

// Cтруктура для CLI-команды
type command struct {
	name string
	agrs []string
}

// Login handler function
func handlerLogin(s *state, cmd command) error {
	// Проверяетм, что в аргументах есть имя
	if len(cmd.agrs) == 0 {
		return fmt.Errorf("Username is required")
	}

	// Проверяем, что пользователь есть в БД
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, cmd.agrs[0])
	if err != nil {
		return fmt.Errorf("user %s does not exist", cmd.agrs[0])
	}

	// Устанавливаем наше имя пользователя
	err = s.cfg.SetUser(cmd.agrs[0])
	if err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		return err
	}

	fmt.Printf("User %s has been logged in\n", s.cfg.CurrentUserName)
	return nil
}

// Register handler fucntion
func handlerRegister(s *state, cmd command) error {
	// Проверяетм, что в аргументах есть имя
		if len(cmd.agrs) == 0 {
		return fmt.Errorf("Username is required")
	}
	
	// Создаем переменные для подставления в CreateUser функцию от sclc
	id := uuid.New()
	now := time.Now()
	name := cmd.agrs[0]
	ctx := context.Background()

	// Формируем структуру из них
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}

	// Проверяем, что такого пользователя нет в БД, если есть - ошибка
	_, err := s.db.GetUser(ctx, name)
	if err == nil {
		return fmt.Errorf("user %s already exists", name)
	}

	// Пытаемся создать пользователя в БД
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// Устанавливаем созданного пользователя в конфиг
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		fmt.Printf("Error setting user: %v\n", err)
		return err
	}

	fmt.Printf("User %s was created\n", name)
	fmt.Println(user)
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
	// Читаем конфиг при запуске
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Инициализируем подключение к БД
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error opening DB: %v\n", err)
		return
	}
	
	dbQueries := database.New(db)

	// Инициализируем состояние и команды
	var s state	
	s.cfg = &cfg
	s.db = dbQueries
	var c commands
	c.handlerMap = make(map[string]func(*state, command) error)

	// Регистрируем команду login
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)

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