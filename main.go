package main

// #cgo CXXFLAGS: -std=c++11
// #cgo LDFLAGS: -lstdc++
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"

	"github.com/AleonDM/MixailOS/core"
	"github.com/AleonDM/MixailOS/ui"
)

var (
	configInstance *core.Config
)

//export GetConfigUsername
func GetConfigUsername() *C.char {
	if configInstance == nil {
		return C.CString("DefaultUser")
	}
	return C.CString(configInstance.Username)
}

//export InitMixailOS
func InitMixailOS() {
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
}

//export StartUI
func StartUI() {
	fmt.Println("Starting MixailOS UI...")
	// Запуск GUI интерфейса
	if configInstance == nil {
		InitMixailOS()
	}
	ui.Run(configInstance)
}

// FreeString освобождает память строки C
//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {
	// Это нужно для корректной работы в Windows
	runtime.LockOSThread()
	
	// Инициализация
	InitMixailOS()
	
	// Запуск UI
	StartUI()
} 