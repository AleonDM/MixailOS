package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AleonDM/MixailOS/core"
	"github.com/AleonDM/MixailOS/ui"
)

var (
	configInstance *core.Config
)

func main() {
	fmt.Println("Initializing MixailOS...")
	
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
	configInstance = core.NewConfig(mixailOSDir)
	if err := configInstance.Load(); err != nil {
		fmt.Println("Загрузка конфигурации по умолчанию")
		configInstance.SetDefault()
		if err := configInstance.Save(); err != nil {
			fmt.Println("Ошибка при сохранении конфигурации:", err)
		}
	}
	
	// Инициализация файловой системы и консоли
	fileSystem := core.NewFileSystem(configInstance)
	console := core.NewConsole(fileSystem, configInstance)
	
	// Запуск GUI интерфейса
	ui.RunUI(configInstance, fileSystem, console)
} 