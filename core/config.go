package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config содержит все настройки MixailOS
type Config struct {
	Username    string `json:"username"`
	Wallpaper   string `json:"wallpaper"`
	RootDir     string `json:"rootDir"`
	CurrentDir  string `json:"currentDir"`
	DefaultApps map[string]string `json:"defaultApps"`
}

// NewConfig создает новый экземпляр конфигурации
func NewConfig(rootDir string) *Config {
	return &Config{
		RootDir:    rootDir,
		CurrentDir: rootDir,
		DefaultApps: map[string]string{
			"browser":  "internal",
			"fileExch": "internal",
			"calc":     "internal",
		},
	}
}

// SetDefault устанавливает значения по умолчанию
func (c *Config) SetDefault() {
	c.Username = "User"
	c.Wallpaper = "default.jpg"
	
	// Создание необходимых директорий
	dirs := []string{
		filepath.Join(c.RootDir, "Documents"),
		filepath.Join(c.RootDir, "Downloads"),
		filepath.Join(c.RootDir, "Pictures"),
		filepath.Join(c.RootDir, "Music"),
		filepath.Join(c.RootDir, "Videos"),
	}
	
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		}
	}
	
	// Создание файла с приветственным сообщением
	welcomeFile := filepath.Join(c.RootDir, "Documents", "welcome.txt")
	welcomeText := "Добро пожаловать в MixailOS!\n\nЭто ваша новая операционная система. Чтобы начать, откройте меню и выберите нужное приложение.\n\nДля вызова консоли нажмите Ctrl+T."
	ioutil.WriteFile(welcomeFile, []byte(welcomeText), 0644)
}

// Load загружает конфигурацию из файла
func (c *Config) Load() error {
	configPath := filepath.Join(c.RootDir, "config.json")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}
	
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, c)
}

// Save сохраняет конфигурацию в файл
func (c *Config) Save() error {
	configPath := filepath.Join(c.RootDir, "config.json")
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(configPath, data, 0644)
}

// ChangeUsername изменяет имя пользователя
func (c *Config) ChangeUsername(newName string) error {
	c.Username = newName
	return c.Save()
}

// ChangeWallpaper изменяет обои рабочего стола
func (c *Config) ChangeWallpaper(wallpaperPath string) error {
	c.Wallpaper = wallpaperPath
	return c.Save()
}

// GetCurrentDir возвращает текущую директорию
func (c *Config) GetCurrentDir() string {
	return c.CurrentDir
}

// SetCurrentDir изменяет текущую директорию
func (c *Config) SetCurrentDir(path string) {
	c.CurrentDir = path
} 