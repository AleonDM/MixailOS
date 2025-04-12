package ui

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/AleonDM/MixailOS/core"
)

type MixailOSUI struct {
	App         fyne.App
	MainWindow  fyne.Window
	Config      *core.Config
	FileSystem  *core.FileSystem
	Console     *core.Console
	
	// Интерфейсные компоненты
	ConsoleOutput *widget.TextGrid
	ConsoleInput  *widget.Entry
	FileList      *widget.List
	CurrentPath   *widget.Label
	UsernameEntry *widget.Entry
}

// RunUI создает и запускает пользовательский интерфейс
func RunUI(config *core.Config, fs *core.FileSystem, console *core.Console) {
	// Инициализация приложения Fyne
	mixailApp := app.New()
	mixailApp.Settings().SetTheme(theme.DarkTheme())
	
	// Создание главного окна
	mainWindow := mixailApp.NewWindow("MixailOS")
	mainWindow.Resize(fyne.NewSize(800, 600))
	
	// Создание экземпляра UI
	ui := &MixailOSUI{
		App:        mixailApp,
		MainWindow: mainWindow,
		Config:     config,
		FileSystem: fs,
		Console:    console,
	}
	
	// Создание интерфейса
	ui.setupUI()
	
	// Отображение и запуск
	mainWindow.ShowAndRun()
}

// setupUI создает все элементы пользовательского интерфейса
func (ui *MixailOSUI) setupUI() {
	// Создание вкладок для разных функций
	tabs := container.NewAppTabs(
		container.NewTabItem("Консоль", ui.createConsoleTab()),
		container.NewTabItem("Файлы", ui.createFileManagerTab()),
		container.NewTabItem("Браузер", ui.createBrowserTab()),
		container.NewTabItem("Калькулятор", ui.createCalculatorTab()),
		container.NewTabItem("Настройки", ui.createSettingsTab()),
	)
	
	// Заголовок с именем пользователя
	userLabel := widget.NewLabel("Пользователь: " + ui.Config.Username)
	
	// Создание менюбара
	menuBar := ui.createMenuBar()
	
	// Размещение всех элементов в главном окне
	ui.MainWindow.SetContent(container.NewBorder(
		menuBar, // top
		nil,     // bottom
		nil,     // left
		nil,     // right
		container.NewBorder(
			userLabel, // top
			nil,       // bottom
			nil,       // left
			nil,       // right
			tabs,
		),
	))
}

// createMenuBar создает верхнее меню приложения
func (ui *MixailOSUI) createMenuBar() fyne.CanvasObject {
	// Пункт меню "Файл"
	fileMenu := fyne.NewMenu("Файл",
		fyne.NewMenuItem("Новый файл", func() {
			ui.createNewFile()
		}),
		fyne.NewMenuItem("Выход", func() {
			ui.MainWindow.Close()
		}),
	)
	
	// Пункт меню "Вид"
	viewMenu := fyne.NewMenu("Вид",
		fyne.NewMenuItem("Светлая тема", func() {
			ui.App.Settings().SetTheme(theme.LightTheme())
		}),
		fyne.NewMenuItem("Тёмная тема", func() {
			ui.App.Settings().SetTheme(theme.DarkTheme())
		}),
	)
	
	// Пункт меню "Справка"
	helpMenu := fyne.NewMenu("Справка",
		fyne.NewMenuItem("О программе", func() {
			ui.showAboutDialog()
		}),
	)
	
	mainMenu := fyne.NewMainMenu(
		fileMenu,
		viewMenu,
		helpMenu,
	)
	
	ui.MainWindow.SetMainMenu(mainMenu)
	
	// Создание тулбара
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			ui.createNewFile()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			ui.showAboutDialog()
		}),
	)
	
	return toolbar
}

// createConsoleTab создает вкладку с консолью
func (ui *MixailOSUI) createConsoleTab() fyne.CanvasObject {
	// Создание поля вывода консоли
	ui.ConsoleOutput = widget.NewTextGrid()
	ui.ConsoleOutput.SetText("Добро пожаловать в консоль MixailOS!\nВведите 'help' для получения списка доступных команд.\n\n")
	
	// Создание поля ввода
	ui.ConsoleInput = widget.NewEntry()
	ui.ConsoleInput.SetPlaceHolder("Введите команду...")
	ui.ConsoleInput.OnSubmitted = func(cmd string) {
		if cmd != "" {
			// Выполнение команды и получение результата
			result := ui.Console.Execute(cmd)
			
			// Обновление вывода консоли
			currentText := ui.ConsoleOutput.Text()
			newText := currentText + ">> " + cmd + "\n" + result + "\n"
			ui.ConsoleOutput.SetText(newText)
			
			// Очистка поля ввода
			ui.ConsoleInput.SetText("")
		}
	}
	
	// Размещение элементов в контейнере
	return container.NewBorder(
		nil, // top
		container.NewBorder(
			nil, // top
			nil, // bottom
			widget.NewLabel(">> "), // left
			nil, // right
			ui.ConsoleInput,
		), // bottom
		nil, // left
		nil, // right
		container.NewScroll(ui.ConsoleOutput),
	)
}

// createFileManagerTab создает вкладку с файловым менеджером
func (ui *MixailOSUI) createFileManagerTab() fyne.CanvasObject {
	// Получаем список файлов
	files, err := ui.FileSystem.ListFiles()
	if err != nil {
		// Если не можем получить список файлов, создаем пустой список
		files = []string{}
	}
	
	// Сохраняем файлы в массив
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file
	}
	
	// Создаем метку с текущим путем
	ui.CurrentPath = widget.NewLabel(ui.Config.GetCurrentDir())
	
	// Кнопка для перехода в родительскую директорию
	upButton := widget.NewButtonWithIcon("Вверх", theme.NavigateBackIcon(), func() {
		parent := filepath.Dir(ui.Config.GetCurrentDir())
		if strings.HasPrefix(parent, ui.Config.RootDir) {
			ui.FileSystem.ChangeDirectory("..")
			ui.refreshFileList()
		}
	})
	
	// Кнопка для создания новой директории
	mkdirButton := widget.NewButtonWithIcon("Новая папка", theme.FolderNewIcon(), func() {
		// Диалог для создания директории
		dirNameEntry := widget.NewEntry()
		dirNameEntry.SetPlaceHolder("Имя новой папки")
		
		dialog.ShowForm("Создать новую папку", "Создать", "Отмена",
			[]*widget.FormItem{
				widget.NewFormItem("Имя:", dirNameEntry),
			},
			func(confirm bool) {
				if confirm && dirNameEntry.Text != "" {
					ui.FileSystem.CreateDirectory(dirNameEntry.Text)
					ui.refreshFileList()
				}
			},
			ui.MainWindow,
		)
	})
	
	// Создаем список файлов с переменной fileNames из внешней области видимости
	ui.FileList = widget.NewList(
		func() int {
			return len(fileNames)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FileIcon()),
				widget.NewLabel("Template Item"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// Проверяем границы массива
			if id < 0 || id >= len(fileNames) {
				return
			}
			
			container := obj.(*fyne.Container)
			label := container.Objects[1].(*widget.Label)
			icon := container.Objects[0].(*widget.Icon)
			
			// Изменяем иконку в зависимости от типа (файл или директория)
			if strings.Contains(fileNames[id], "(dir)") {
				icon.SetResource(theme.FolderIcon())
			} else {
				icon.SetResource(theme.FileIcon())
			}
			
			label.SetText(fileNames[id])
		},
	)
	
	ui.FileList.OnSelected = func(id widget.ListItemID) {
		// Проверяем границы массива
		if id < 0 || id >= len(fileNames) {
			return
		}
		
		fileName := fileNames[id]
		// Извлекаем имя файла без типа
		parts := strings.Split(fileName, " (")
		name := parts[0]
		
		// Проверяем, директория ли это
		if strings.Contains(fileName, "(dir)") {
			ui.FileSystem.ChangeDirectory(name)
			ui.refreshFileList()
		} else if strings.HasSuffix(name, ".txt") {
			// Читаем текстовый файл
			content, err := ui.FileSystem.ReadTextFile(name)
			if err != nil {
				dialog.ShowError(err, ui.MainWindow)
				return
			}
			
			// Показываем диалог с содержимым
			textViewer := widget.NewMultiLineEntry()
			textViewer.SetText(content)
			textViewer.Disable() // Только для чтения
			
			dialog.ShowCustom("Файл: "+name, "Закрыть", container.NewScroll(textViewer), ui.MainWindow)
		}
	}
	
	// Toolbar для файлового менеджера
	fileToolbar := container.NewHBox(
		upButton,
		mkdirButton,
		widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
			// Проверяем, выбран ли файл
			if ui.FileList.Selected() < 0 || ui.FileList.Selected() >= len(fileNames) {
				dialog.ShowInformation("Внимание", "Выберите файл для удаления", ui.MainWindow)
				return
			}
			
			fileName := fileNames[ui.FileList.Selected()]
			parts := strings.Split(fileName, " (")
			name := parts[0]
			
			dialog.ShowConfirm("Подтверждение", "Вы уверены, что хотите удалить "+name+"?",
				func(confirm bool) {
					if confirm {
						ui.FileSystem.DeleteFile(name)
						ui.refreshFileList()
					}
				},
				ui.MainWindow,
			)
		}),
		widget.NewButtonWithIcon("Обновить", theme.ViewRefreshIcon(), func() {
			ui.refreshFileList()
		}),
	)
	
	// Размещение элементов в контейнере
	return container.NewBorder(
		container.NewVBox(
			ui.CurrentPath,
			fileToolbar,
		), // top
		nil, // bottom
		nil, // left
		nil, // right
		ui.FileList,
	)
}

// createBrowserTab создает вкладку с браузером
func (ui *MixailOSUI) createBrowserTab() fyne.CanvasObject {
	// Поле для ввода URL
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Введите URL...")
	
	// Кнопка для перехода по URL
	goButton := widget.NewButtonWithIcon("Перейти", theme.NavigateNextIcon(), func() {
		// Здесь будет код для загрузки страницы
		// В текущей версии просто показываем заглушку
		content := fmt.Sprintf("Просмотр страницы: %s\n\n--- Содержимое страницы ---\nЭто упрощенная версия браузера MixailOS.\nЗдесь вы увидите содержимое загруженной страницы.", urlEntry.Text)
		
		browserContent.SetText(content)
	})
	
	// Поле для отображения содержимого страницы
	browserContent := widget.NewMultiLineEntry()
	browserContent.SetText("Добро пожаловать в браузер MixailOS!\nВведите URL в поле выше и нажмите 'Перейти' для начала работы.")
	browserContent.Disable() // Только для чтения
	
	// URL bar
	urlBar := container.NewBorder(
		nil, // top
		nil, // bottom
		widget.NewLabel("URL:"), // left
		goButton, // right
		urlEntry,
	)
	
	// Размещение элементов в контейнере
	return container.NewBorder(
		urlBar, // top
		nil, // bottom
		nil, // left
		nil, // right
		container.NewScroll(browserContent),
	)
}

// createCalculatorTab создает вкладку с калькулятором
func (ui *MixailOSUI) createCalculatorTab() fyne.CanvasObject {
	// Поле для отображения результата
	result := widget.NewEntry()
	result.SetPlaceHolder("0")
	
	// Поле для истории вычислений
	history := widget.NewMultiLineEntry()
	history.SetText("Калькулятор MixailOS\nПоддерживаемые операции: +, -, *, /\n\n")
	history.Disable() // Только для чтения
	
	// Создаем кнопки цифр и операций
	var currentInput string
	
	// Функция для добавления символа к вводу
	appendToInput := func(s string) {
		currentInput += s
		result.SetText(currentInput)
	}
	
	// Функция для вычисления результата
	calculate := func() {
		// Если входная строка пуста, ничего не делаем
		if currentInput == "" {
			return
		}
		
		// Разбиваем входную строку на части (числа и операции)
		parts := strings.Fields(currentInput)
		if len(parts) < 3 || len(parts)%2 == 0 {
			// Если формат не соответствует A op B op C...
			historyText := history.Text
			historyText += currentInput + " = Ошибка формата\n"
			history.SetText(historyText)
			currentInput = ""
			result.SetText("Ошибка")
			return
		}
		
		// Выполняем операции последовательно слева направо
		// Сначала обрабатываем умножение и деление
		var val float64
		var err error
		
		// Первое число
		val, err = parseFloat(parts[0])
		if err != nil {
			historyText := history.Text
			historyText += currentInput + " = Ошибка: " + err.Error() + "\n"
			history.SetText(historyText)
			currentInput = ""
			result.SetText("Ошибка")
			return
		}
		
		// Первый проход: умножение и деление
		for i := 1; i < len(parts); i += 2 {
			op := parts[i]
			if op != "*" && op != "/" {
				continue
			}
			
			// Получаем следующее число
			next, err := parseFloat(parts[i+1])
			if err != nil {
				historyText := history.Text
				historyText += currentInput + " = Ошибка: " + err.Error() + "\n"
				history.SetText(historyText)
				currentInput = ""
				result.SetText("Ошибка")
				return
			}
			
			// Выполняем операцию
			switch op {
			case "*":
				val *= next
				// Обновляем части, чтобы пропустить этот результат при последующих операциях
				parts[i+1] = fmt.Sprintf("%g", val)
				parts[i] = " " // пометка для пропуска
				parts[i-1] = " "
			case "/":
				if next == 0 {
					historyText := history.Text
					historyText += currentInput + " = Ошибка: деление на ноль\n"
					history.SetText(historyText)
					currentInput = ""
					result.SetText("Ошибка")
					return
				}
				val /= next
				// Обновляем части, чтобы пропустить этот результат при последующих операциях
				parts[i+1] = fmt.Sprintf("%g", val)
				parts[i] = " " // пометка для пропуска
				parts[i-1] = " "
			}
		}
		
		// Второй проход: сложение и вычитание
		foundValue := false
		val = 0
		
		for i := 0; i < len(parts); i++ {
			if parts[i] == " " {
				continue
			}
			
			if !foundValue {
				val, err = parseFloat(parts[i])
				if err == nil {
					foundValue = true
				}
				continue
			}
			
			if parts[i] == "+" || parts[i] == "-" {
				op := parts[i]
				if i+1 < len(parts) && parts[i+1] != " " {
					next, err := parseFloat(parts[i+1])
					if err != nil {
						continue
					}
					
					switch op {
					case "+":
						val += next
					case "-":
						val -= next
					}
				}
			}
		}
		
		// Форматируем результат и обновляем историю
		resultStr := fmt.Sprintf("%g", val)
		historyText := history.Text
		historyText += currentInput + " = " + resultStr + "\n"
		history.SetText(historyText)
		
		// Обновляем поле результата
		result.SetText(resultStr)
		currentInput = resultStr
	}
	
	// Кнопки для цифр
	digits := container.NewGridWithColumns(3,
		widget.NewButton("7", func() { appendToInput("7") }),
		widget.NewButton("8", func() { appendToInput("8") }),
		widget.NewButton("9", func() { appendToInput("9") }),
		widget.NewButton("4", func() { appendToInput("4") }),
		widget.NewButton("5", func() { appendToInput("5") }),
		widget.NewButton("6", func() { appendToInput("6") }),
		widget.NewButton("1", func() { appendToInput("1") }),
		widget.NewButton("2", func() { appendToInput("2") }),
		widget.NewButton("3", func() { appendToInput("3") }),
		widget.NewButton("0", func() { appendToInput("0") }),
		widget.NewButton(".", func() { appendToInput(".") }),
		widget.NewButton("=", calculate),
	)
	
	// Кнопки для операций
	operations := container.NewVBox(
		widget.NewButton("+", func() { appendToInput(" + ") }),
		widget.NewButton("-", func() { appendToInput(" - ") }),
		widget.NewButton("*", func() { appendToInput(" * ") }),
		widget.NewButton("/", func() { appendToInput(" / ") }),
		widget.NewButton("C", func() {
			currentInput = ""
			result.SetText("0")
		}),
	)
	
	// Размещение элементов в контейнере
	calcLayout := container.NewBorder(
		container.NewVBox(
			result,
			container.NewScroll(history),
		), // top
		nil, // bottom
		nil, // left
		operations, // right
		digits,
	)
	
	return calcLayout
}

// parseFloat преобразует строку в число с плавающей точкой
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// createSettingsTab создает вкладку с настройками
func (ui *MixailOSUI) createSettingsTab() fyne.CanvasObject {
	// Поле для изменения имени пользователя
	ui.UsernameEntry = widget.NewEntry()
	ui.UsernameEntry.SetText(ui.Config.Username)
	
	usernameForm := widget.NewForm(
		widget.NewFormItem("Имя пользователя:", ui.UsernameEntry),
	)
	
	saveUsernameButton := widget.NewButton("Сохранить имя пользователя", func() {
		if ui.UsernameEntry.Text != "" {
			ui.Config.ChangeUsername(ui.UsernameEntry.Text)
			dialog.ShowInformation("Успех", "Имя пользователя успешно изменено", ui.MainWindow)
		}
	})
	
	// Кнопка для изменения обоев
	changeWallpaperButton := widget.NewButton("Изменить обои", func() {
		// Диалог выбора файла
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, ui.MainWindow)
				return
			}
			if uri == nil {
				return
			}
			
			// Получаем путь к файлу
			path := uri.URI().Path()
			
			// Изменяем обои
			ui.Config.ChangeWallpaper(path)
			dialog.ShowInformation("Успех", "Обои успешно изменены", ui.MainWindow)
		}, ui.MainWindow)
	})
	
	// Размещение элементов в контейнере
	return container.NewVBox(
		widget.NewLabel("Настройки MixailOS"),
		usernameForm,
		saveUsernameButton,
		widget.NewSeparator(),
		changeWallpaperButton,
	)
}

// createNewFile создает новый текстовый файл
func (ui *MixailOSUI) createNewFile() {
	// Создаем диалог для создания файла
	filenameEntry := widget.NewEntry()
	filenameEntry.SetPlaceHolder("имя_файла.txt")
	
	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Содержимое файла...")
	
	dialog.ShowForm("Создать новый файл", "Создать", "Отмена",
		[]*widget.FormItem{
			widget.NewFormItem("Имя файла:", filenameEntry),
			widget.NewFormItem("Содержимое:", contentEntry),
		},
		func(confirm bool) {
			if confirm && filenameEntry.Text != "" {
				// Создаем файл
				ui.FileSystem.CreateTextFile(filenameEntry.Text, contentEntry.Text)
				
				// Обновляем список файлов, если находимся в файловом менеджере
				ui.refreshFileList()
			}
		},
		ui.MainWindow,
	)
}

// showAboutDialog показывает диалог "О программе"
func (ui *MixailOSUI) showAboutDialog() {
	aboutText := "MixailOS v1.0\n\nЭмулятор операционной системы на Go.\nРазработано с использованием Fyne.io\n\n© 2024"
	dialog.ShowInformation("О программе", aboutText, ui.MainWindow)
}

// refreshFileList обновляет список файлов в UI
func (ui *MixailOSUI) refreshFileList() {
	// Обновляем текст текущего пути
	ui.CurrentPath.SetText(ui.Config.GetCurrentDir())
	
	// Получаем список файлов
	files, err := ui.FileSystem.ListFiles()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Ошибка при получении списка файлов: %v", err), ui.MainWindow)
		return
	}
	
	// Сохраняем файлы в новый массив
	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file
	}
	
	// Создаем новый список файлов
	newFileList := widget.NewList(
		func() int {
			return len(fileNames)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FileIcon()),
				widget.NewLabel("Template Item"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// Проверяем границы массива
			if id < 0 || id >= len(fileNames) {
				return
			}
			
			container := obj.(*fyne.Container)
			label := container.Objects[1].(*widget.Label)
			icon := container.Objects[0].(*widget.Icon)
			
			// Изменяем иконку в зависимости от типа (файл или директория)
			if strings.Contains(fileNames[id], "(dir)") {
				icon.SetResource(theme.FolderIcon())
			} else {
				icon.SetResource(theme.FileIcon())
			}
			
			label.SetText(fileNames[id])
		},
	)
	
	newFileList.OnSelected = func(id widget.ListItemID) {
		// Проверяем границы массива
		if id < 0 || id >= len(fileNames) {
			return
		}
		
		fileName := fileNames[id]
		// Извлекаем имя файла без типа
		parts := strings.Split(fileName, " (")
		name := parts[0]
		
		// Проверяем, директория ли это
		if strings.Contains(fileName, "(dir)") {
			ui.FileSystem.ChangeDirectory(name)
			ui.refreshFileList()
		} else if strings.HasSuffix(name, ".txt") {
			// Читаем текстовый файл
			content, err := ui.FileSystem.ReadTextFile(name)
			if err != nil {
				dialog.ShowError(err, ui.MainWindow)
				return
			}
			
			// Показываем диалог с содержимым
			textViewer := widget.NewMultiLineEntry()
			textViewer.SetText(content)
			textViewer.Disable() // Только для чтения
			
			dialog.ShowCustom("Файл: "+name, "Закрыть", container.NewScroll(textViewer), ui.MainWindow)
		}
	}
	
	// Находим контейнер, содержащий старый список файлов
	tabs := ui.MainTabs.Objects[0].(*container.AppTabs)
	fileManagerTab := tabs.Items[1].Content.(*fyne.Container)
	fileListContainer := fileManagerTab.Objects[0].(*fyne.Container)
	
	// Удаляем старый список файлов
	fileListContainer.Objects[0] = newFileList
	ui.FileList = newFileList
	
	// Обновляем UI
	fileListContainer.Refresh()
} 