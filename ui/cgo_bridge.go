package ui

//go:generate go run golang.org/x/tools/cmd/cgo -exportheader ui_bridge.h

// #cgo CXXFLAGS: -std=c++11
// #cgo LDFLAGS: -lstdc++
// #include <stdlib.h>
// #include "ui_bridge.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/AleonDM/MixailOS/core"
)

var (
	globalConfig     *core.Config
	globalFileSystem *core.FileSystem
	globalConsole    *core.Console
)

// Initialize инициализирует интерфейс и компоненты
func Initialize(config *core.Config) {
	fmt.Println("Initializing UI components...")
	globalConfig = config
	globalFileSystem = core.NewFileSystem(config)
	globalConsole = core.NewConsole(globalFileSystem, config)
}

//export GoGetUsername
func GoGetUsername() *C.char {
	if globalConfig == nil {
		return C.CString("DefaultUser")
	}
	return C.CString(globalConfig.Username)
}

//export GoSetUsername
func GoSetUsername(cUsername *C.char) {
	if globalConfig == nil {
		fmt.Println("Error: Config is not initialized!")
		return
	}
	username := C.GoString(cUsername)
	globalConfig.ChangeUsername(username)
}

//export GoGetCurrentDirectory
func GoGetCurrentDirectory() *C.char {
	if globalConfig == nil {
		return C.CString(".")
	}
	return C.CString(globalConfig.GetCurrentDir())
}

//export GoExecuteConsoleCommand
func GoExecuteConsoleCommand(cCmd *C.char) *C.char {
	if globalConsole == nil {
		return C.CString("Error: Console is not initialized!")
	}
	cmd := C.GoString(cCmd)
	result := globalConsole.Execute(cmd)
	return C.CString(result)
}

//export GoGetFileList
func GoGetFileList() *C.char {
	if globalFileSystem == nil {
		return C.CString("Error: FileSystem is not initialized!")
	}
	
	files, err := globalFileSystem.ListFiles()
	if err != nil {
		return C.CString("Ошибка: " + err.Error())
	}
	
	result := ""
	for i, file := range files {
		if i > 0 {
			result += "|"
		}
		result += file
	}
	
	return C.CString(result)
}

//export GoChangeWallpaper
func GoChangeWallpaper(cPath *C.char) {
	if globalConfig == nil {
		fmt.Println("Error: Config is not initialized!")
		return
	}
	path := C.GoString(cPath)
	globalConfig.ChangeWallpaper(path)
	globalConfig.Save()
}

//export GoGetWallpaperPath
func GoGetWallpaperPath() *C.char {
	if globalConfig == nil {
		return C.CString("default.jpg")
	}
	return C.CString(globalConfig.Wallpaper)
}

//export GoCreateTextFile
func GoCreateTextFile(cName *C.char, cContent *C.char) *C.char {
	if globalFileSystem == nil {
		return C.CString("Error: FileSystem is not initialized!")
	}
	name := C.GoString(cName)
	content := C.GoString(cContent)
	
	err := globalFileSystem.CreateTextFile(name, content)
	if err != nil {
		return C.CString("Ошибка: " + err.Error())
	}
	
	return C.CString("Файл успешно создан")
}

//export GoReadTextFile
func GoReadTextFile(cName *C.char) *C.char {
	if globalFileSystem == nil {
		return C.CString("Error: FileSystem is not initialized!")
	}
	name := C.GoString(cName)
	
	content, err := globalFileSystem.ReadTextFile(name)
	if err != nil {
		return C.CString("Ошибка: " + err.Error())
	}
	
	return C.CString(content)
}

//export RunUI
func RunUI() {
	fmt.Println("RunUI called from Go!")
	if globalConfig == nil {
		fmt.Println("Warning: config is nil, initializing...")
		// У нас есть экспортированная функция InitMixailOS, но мы не можем её вызвать напрямую
		// поэтому UI не запустится корректно без предварительной инициализации
	}
	C.RunUI()
}

// Run запускает пользовательский интерфейс
func Run(config *core.Config) {
	fmt.Println("Starting UI from Go...")
	Initialize(config)
	RunUI()
}

//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
} 