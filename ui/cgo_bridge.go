package ui

// #cgo CXXFLAGS: -std=c++11
// #cgo LDFLAGS: -lstdc++
// #include <stdlib.h>
// #include "ui_bridge.h"
import "C"
import (
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
	globalConfig = config
	globalFileSystem = core.NewFileSystem(config)
	globalConsole = core.NewConsole(globalFileSystem, config)
}

//export GoGetUsername
func GoGetUsername() *C.char {
	return C.CString(globalConfig.Username)
}

//export GoSetUsername
func GoSetUsername(cUsername *C.char) {
	username := C.GoString(cUsername)
	globalConfig.ChangeUsername(username)
}

//export GoGetCurrentDirectory
func GoGetCurrentDirectory() *C.char {
	return C.CString(globalConfig.GetCurrentDir())
}

//export GoExecuteConsoleCommand
func GoExecuteConsoleCommand(cCmd *C.char) *C.char {
	cmd := C.GoString(cCmd)
	result := globalConsole.Execute(cmd)
	return C.CString(result)
}

//export GoGetFileList
func GoGetFileList() *C.char {
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
	path := C.GoString(cPath)
	globalConfig.ChangeWallpaper(path)
	globalConfig.Save()
}

//export GoGetWallpaperPath
func GoGetWallpaperPath() *C.char {
	return C.CString(globalConfig.Wallpaper)
}

//export GoCreateTextFile
func GoCreateTextFile(cName *C.char, cContent *C.char) *C.char {
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
	name := C.GoString(cName)
	
	content, err := globalFileSystem.ReadTextFile(name)
	if err != nil {
		return C.CString("Ошибка: " + err.Error())
	}
	
	return C.CString(content)
}

// Run запускает пользовательский интерфейс
func Run(config *core.Config) {
	Initialize(config)
	C.RunUI()
}

// FreeString освобождает память строки C
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
} 