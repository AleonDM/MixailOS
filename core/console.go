package core

import (
	"fmt"
	"strings"
	"time"
)

// Console представляет интерфейс консоли MixailOS
type Console struct {
	FileSystem *FileSystem
	Config     *Config
	History    []string
}

// NewConsole создает новый экземпляр консоли
func NewConsole(fs *FileSystem, config *Config) *Console {
	return &Console{
		FileSystem: fs,
		Config:     config,
		History:    []string{},
	}
}

// Execute выполняет команду консоли и возвращает результат
func (c *Console) Execute(cmd string) string {
	c.History = append(c.History, cmd)
	
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	
	switch parts[0] {
	case "help":
		return c.HelpCommand()
	case "info":
		return c.InfoCommand()
	case "cls":
		return "clear"
	case "txt":
		if len(parts) < 2 {
			return "Использование: txt [read|write|list] [параметры]"
		}
		return c.TextCommand(parts[1], parts[2:])
	case "cd":
		if len(parts) < 2 {
			return fmt.Sprintf("Текущая директория: %s", c.Config.CurrentDir)
		}
		return c.CdCommand(parts[1])
	case "ls":
		return c.LsCommand()
	case "mkdir":
		if len(parts) < 2 {
			return "Использование: mkdir <имя_директории>"
		}
		return c.MkdirCommand(parts[1])
	case "rm":
		if len(parts) < 2 {
			return "Использование: rm <имя_файла>"
		}
		return c.RmCommand(parts[1])
	case "cp":
		if len(parts) < 3 {
			return "Использование: cp <исходный_файл> <файл_назначения>"
		}
		return c.CpCommand(parts[1], parts[2])
	case "echo":
		return strings.Join(parts[1:], " ")
	case "date":
		return time.Now().Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("Неизвестная команда: %s. Введите 'help' для получения списка команд.", parts[0])
	}
}

// HelpCommand возвращает справку по командам
func (c *Console) HelpCommand() string {
	return `Доступные команды:
help - показать список команд
info - показать информацию о системе
cls - очистить экран консоли
txt - работа с текстовыми файлами:
  - txt read <имя_файла> - чтение файла
  - txt write <имя_файла> <содержимое> - запись в файл
  - txt list - список текстовых файлов
cd <путь> - изменить текущую директорию
ls - показать содержимое текущей директории
mkdir <имя> - создать новую директорию
rm <имя> - удалить файл
cp <источник> <назначение> - копировать файл
echo <текст> - вывести текст
date - показать текущую дату и время`
}

// InfoCommand возвращает информацию о системе
func (c *Console) InfoCommand() string {
	return fmt.Sprintf(`Информация о системе MixailOS:
Пользователь: %s
Рабочая директория: %s
Текущая директория: %s
Версия: 1.0.0
Дата запуска: %s`, 
		c.Config.Username,
		c.Config.RootDir,
		c.Config.CurrentDir,
		time.Now().Format("2006-01-02 15:04:05"))
}

// TextCommand обрабатывает команды для работы с текстовыми файлами
func (c *Console) TextCommand(action string, args []string) string {
	switch action {
	case "read":
		if len(args) < 1 {
			return "Использование: txt read <имя_файла>"
		}
		content, err := c.FileSystem.ReadTextFile(args[0])
		if err != nil {
			return fmt.Sprintf("Ошибка при чтении файла: %v", err)
		}
		return fmt.Sprintf("Содержимое файла %s:\n%s", args[0], content)
		
	case "write":
		if len(args) < 2 {
			return "Использование: txt write <имя_файла> <содержимое>"
		}
		content := strings.Join(args[1:], " ")
		err := c.FileSystem.CreateTextFile(args[0], content)
		if err != nil {
			return fmt.Sprintf("Ошибка при записи файла: %v", err)
		}
		return fmt.Sprintf("Файл %s успешно создан", args[0])
		
	case "list":
		files, err := c.FileSystem.ListFiles()
		if err != nil {
			return fmt.Sprintf("Ошибка при получении списка файлов: %v", err)
		}
		
		var textFiles []string
		for _, file := range files {
			if strings.Contains(file, ".txt") {
				textFiles = append(textFiles, file)
			}
		}
		
		if len(textFiles) == 0 {
			return "Текстовые файлы не найдены"
		}
		
		return fmt.Sprintf("Текстовые файлы:\n%s", strings.Join(textFiles, "\n"))
		
	default:
		return fmt.Sprintf("Неизвестное действие для txt: %s", action)
	}
}

// CdCommand изменяет текущую директорию
func (c *Console) CdCommand(path string) string {
	err := c.FileSystem.ChangeDirectory(path)
	if err != nil {
		return fmt.Sprintf("Ошибка при изменении директории: %v", err)
	}
	return fmt.Sprintf("Текущая директория: %s", c.Config.CurrentDir)
}

// LsCommand показывает содержимое текущей директории
func (c *Console) LsCommand() string {
	files, err := c.FileSystem.ListFiles()
	if err != nil {
		return fmt.Sprintf("Ошибка при получении списка файлов: %v", err)
	}
	
	if len(files) == 0 {
		return "Директория пуста"
	}
	
	return fmt.Sprintf("Содержимое директории %s:\n%s", c.Config.CurrentDir, strings.Join(files, "\n"))
}

// MkdirCommand создает новую директорию
func (c *Console) MkdirCommand(name string) string {
	err := c.FileSystem.CreateDirectory(name)
	if err != nil {
		return fmt.Sprintf("Ошибка при создании директории: %v", err)
	}
	return fmt.Sprintf("Директория %s успешно создана", name)
}

// RmCommand удаляет файл
func (c *Console) RmCommand(name string) string {
	err := c.FileSystem.DeleteFile(name)
	if err != nil {
		return fmt.Sprintf("Ошибка при удалении файла: %v", err)
	}
	return fmt.Sprintf("Файл %s успешно удален", name)
}

// CpCommand копирует файл
func (c *Console) CpCommand(src, dst string) string {
	err := c.FileSystem.CopyFile(src, dst)
	if err != nil {
		return fmt.Sprintf("Ошибка при копировании файла: %v", err)
	}
	return fmt.Sprintf("Файл %s успешно скопирован в %s", src, dst)
} 