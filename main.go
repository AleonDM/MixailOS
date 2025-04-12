package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AleonDM/MixailOS/core"
	"github.com/AleonDM/MixailOS/ui"
)

func main() {
	// Инициализация рабочей директории
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Ошибка при получении домашней директории:", err)
		os.Exit(1)
	}

	mixailOSDir := filepath.Join(homeDir, "MixailOS")
	if _, err := os.Stat(mixailOSDir); os.IsNotExist(err) {
		if err := os.Mkdir(mixailOSDir, 0755); err != nil {
			fmt.Println("Ошибка при создании директории MixailOS:", err)
			os.Exit(1)
		}
	}

	// Инициализация системных настроек
	config := core.NewConfig(mixailOSDir)
	if err := config.Load(); err != nil {
		fmt.Println("Загрузка конфигурации по умолчанию")
		config.SetDefault()
		if err := config.Save(); err != nil {
			fmt.Println("Ошибка при сохранении конфигурации:", err)
		}
	}

	// Запуск GUI интерфейса
	ui.Run(config)
} 