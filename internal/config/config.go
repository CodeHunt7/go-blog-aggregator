package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	// 1. Получаем путь к домашней директории
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// 2. Собираем полный путь к файлу
	fullPath := filepath.Join(home, configFileName)

	return fullPath, nil
}

func Read() (Config, error) {
	// 1. Получаем путь к файлу
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	// 2. Читаем файл
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	// 3. Декодируем JSON в структуру
	var cfg Config
	err = json.Unmarshal(fileData, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error { 
	// 1. Записываем имя в структуру
	cfg.CurrentUserName = userName

	// 2. Пуляем структуру в JSON
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// 3. Получаем путь к домашней директории и собираем путь к файлу
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// 4. записываем JSON в файл
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	
	return nil
}