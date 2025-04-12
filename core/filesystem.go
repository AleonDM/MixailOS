package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem представляет файловую систему MixailOS
type FileSystem struct {
	Config *Config
}

// NewFileSystem создает новый экземпляр файловой системы
func NewFileSystem(config *Config) *FileSystem {
	return &FileSystem{
		Config: config,
	}
}

// ListFiles возвращает список файлов и директорий в текущей директории
func (fs *FileSystem) ListFiles() ([]string, error) {
	files, err := ioutil.ReadDir(fs.Config.CurrentDir)
	if err != nil {
		return nil, err
	}
	
	var fileNames []string
	for _, file := range files {
		fileType := "file"
		if file.IsDir() {
			fileType = "dir"
		}
		fileNames = append(fileNames, fmt.Sprintf("%s (%s)", file.Name(), fileType))
	}
	
	return fileNames, nil
}

// ChangeDirectory изменяет текущую директорию
func (fs *FileSystem) ChangeDirectory(path string) error {
	// Обработка специальных символов
	if path == ".." {
		parent := filepath.Dir(fs.Config.CurrentDir)
		// Не позволяем выйти за пределы рабочей директории MixailOS
		if strings.HasPrefix(parent, fs.Config.RootDir) {
			fs.Config.CurrentDir = parent
		}
		return nil
	}
	
	// Если путь относительный, добавляем текущую директорию
	if !filepath.IsAbs(path) {
		path = filepath.Join(fs.Config.CurrentDir, path)
	}
	
	// Не позволяем выйти за пределы рабочей директории MixailOS
	if !strings.HasPrefix(path, fs.Config.RootDir) {
		return fmt.Errorf("недопустимый путь: %s", path)
	}
	
	// Проверяем, существует ли директория
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s не является директорией", path)
	}
	
	fs.Config.CurrentDir = path
	return nil
}

// CreateTextFile создает текстовый файл с указанным содержимым
func (fs *FileSystem) CreateTextFile(name string, content string) error {
	// Добавляем расширение .txt, если его нет
	if !strings.HasSuffix(name, ".txt") {
		name = name + ".txt"
	}
	
	path := filepath.Join(fs.Config.CurrentDir, name)
	return ioutil.WriteFile(path, []byte(content), 0644)
}

// ReadTextFile читает текстовый файл
func (fs *FileSystem) ReadTextFile(name string) (string, error) {
	// Добавляем расширение .txt, если его нет
	if !strings.HasSuffix(name, ".txt") {
		name = name + ".txt"
	}
	
	path := filepath.Join(fs.Config.CurrentDir, name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// DeleteFile удаляет файл
func (fs *FileSystem) DeleteFile(name string) error {
	path := filepath.Join(fs.Config.CurrentDir, name)
	
	// Проверяем, существует ли файл
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	
	return os.Remove(path)
}

// CreateDirectory создает новую директорию
func (fs *FileSystem) CreateDirectory(name string) error {
	path := filepath.Join(fs.Config.CurrentDir, name)
	return os.Mkdir(path, 0755)
}

// CopyFile копирует файл
func (fs *FileSystem) CopyFile(src, dst string) error {
	srcPath := filepath.Join(fs.Config.CurrentDir, src)
	dstPath := filepath.Join(fs.Config.CurrentDir, dst)
	
	input, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(dstPath, input, 0644)
} 