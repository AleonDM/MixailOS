#include <FL/Fl.H>
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Box.H>
#include <FL/Fl_Button.H>
#include <FL/Fl_Input.H>
#include <FL/Fl_Text_Display.H>
#include <FL/Fl_Text_Buffer.H>
#include <FL/Fl_Menu_Bar.H>
#include <FL/Fl_File_Chooser.H>
#include <FL/Fl_Tabs.H>
#include <FL/Fl_Group.H>
#include <FL/Fl_PNG_Image.H>
#include <FL/Fl_JPEG_Image.H>
#include <FL/Fl_Shared_Image.H>
#include <string>
#include <vector>
#include <cstdlib>
#include <cstring>
#include <iostream>
#include <sstream>

#include "ui_bridge.h"

// Глобальные переменные для главных компонентов UI
Fl_Double_Window* mainWindow = nullptr;
Fl_Text_Buffer* consoleBuffer = nullptr;
Fl_Text_Display* consoleDisplay = nullptr;
Fl_Input* consoleInput = nullptr;
Fl_Box* wallpaperBox = nullptr;
Fl_Input* usernameInput = nullptr;
Fl_Text_Buffer* browserBuffer = nullptr;
Fl_Text_Display* browserDisplay = nullptr;
Fl_Input* browserUrlInput = nullptr;
Fl_Input* calculatorInput = nullptr;
Fl_Text_Display* calculatorDisplay = nullptr;
Fl_Text_Buffer* calculatorBuffer = nullptr;

// Функции для освобождения памяти строк C
void freeString(char* str) {
    if (str) {
        free(str);
    }
}

// Обновление обоев
void updateWallpaper() {
    char* wallpaperPath = GoGetWallpaperPath();
    
    if (wallpaperPath && strlen(wallpaperPath) > 0) {
        // Попробуем загрузить изображение
        Fl_Shared_Image *img = nullptr;
        
        if (strstr(wallpaperPath, ".jpg") || strstr(wallpaperPath, ".jpeg")) {
            img = Fl_Shared_Image::get(wallpaperPath);
        } else if (strstr(wallpaperPath, ".png")) {
            img = Fl_Shared_Image::get(wallpaperPath);
        }
        
        if (img && img->w() > 0 && img->h() > 0) {
            // Масштабируем изображение под размер фона
            Fl_Image *scaled = img->copy(mainWindow->w(), mainWindow->h());
            wallpaperBox->image(scaled);
            wallpaperBox->redraw();
        }
    }
    
    freeString(wallpaperPath);
}

// Обработчик для ввода команд в консоль
void consoleInputCallback(Fl_Widget* w, void*) {
    Fl_Input* input = static_cast<Fl_Input*>(w);
    const char* text = input->value();
    
    if (text && strlen(text) > 0) {
        // Отправляем команду в Go и получаем результат
        char* result = GoExecuteConsoleCommand(const_cast<char*>(text));
        
        // Показываем команду и результат в консоли
        consoleBuffer->append(">> ");
        consoleBuffer->append(text);
        consoleBuffer->append("\n");
        
        // Если команда "cls", очищаем консоль
        if (strcmp(text, "cls") == 0) {
            consoleBuffer->text("");
        } else {
            consoleBuffer->append(result);
            consoleBuffer->append("\n");
        }
        
        freeString(result);
        input->value(""); // Очищаем поле ввода
        
        // Прокручиваем до конца
        consoleDisplay->scroll(consoleBuffer->length(), 0);
    }
}

// Обработчик для изменения имени пользователя
void usernameCallback(Fl_Widget* w, void*) {
    Fl_Input* input = static_cast<Fl_Input*>(w);
    const char* username = input->value();
    
    if (username && strlen(username) > 0) {
        GoSetUsername(const_cast<char*>(username));
        
        // Обновляем заголовок окна
        std::string title = "MixailOS - " + std::string(username);
        mainWindow->label(title.c_str());
    }
}

// Обработчик для кнопки изменения обоев
void changeWallpaperCallback(Fl_Widget*, void*) {
    Fl_File_Chooser chooser(".", "*.{jpg,jpeg,png}", Fl_File_Chooser::SINGLE, "Выберите изображение для обоев");
    chooser.show();
    
    while (chooser.shown()) {
        Fl::wait();
    }
    
    if (chooser.value()) {
        GoChangeWallpaper(const_cast<char*>(chooser.value()));
        updateWallpaper();
    }
}

// Вычисление в калькуляторе
void calculateCallback(Fl_Widget*, void*) {
    const char* expr = calculatorInput->value();
    
    // Простой калькулятор, который обрабатывает базовые операции
    // В реальном приложении здесь должен быть более сложный парсер выражений
    double result = 0;
    char op = '+';
    
    // Очень простой парсер выражений
    std::stringstream ss(expr);
    double number;
    
    ss >> result;
    
    while (ss >> op >> number) {
        switch (op) {
            case '+':
                result += number;
                break;
            case '-':
                result -= number;
                break;
            case '*':
                result *= number;
                break;
            case '/':
                if (number != 0) {
                    result /= number;
                } else {
                    calculatorBuffer->text("Ошибка: деление на ноль\n");
                    return;
                }
                break;
            default:
                calculatorBuffer->text("Неподдерживаемая операция\n");
                return;
        }
    }
    
    // Выводим результат
    char resultStr[100];
    snprintf(resultStr, sizeof(resultStr), "%s = %.2f\n", expr, result);
    calculatorBuffer->append(resultStr);
    calculatorDisplay->scroll(calculatorBuffer->length(), 0);
}

// Обработчик для браузера
void browserGoCallback(Fl_Widget*, void*) {
    const char* url = browserUrlInput->value();
    
    if (url && strlen(url) > 0) {
        // Базовая реализация браузера - просто отображает URL и заглушку
        std::string content = "Просмотр страницы: ";
        content += url;
        content += "\n\n--- Содержимое страницы ---\n";
        content += "Это упрощенная версия браузера MixailOS.\n";
        content += "Здесь вы увидите содержимое загруженной страницы.\n";
        
        browserBuffer->text(content.c_str());
    }
}

// Создание вкладки для консоли
Fl_Group* createConsoleTab(int x, int y, int w, int h) {
    Fl_Group* consoleTab = new Fl_Group(x, y, w, h, "Консоль");
    
    consoleBuffer = new Fl_Text_Buffer();
    consoleDisplay = new Fl_Text_Display(x + 10, y + 10, w - 20, h - 50);
    consoleDisplay->buffer(consoleBuffer);
    consoleDisplay->textfont(FL_COURIER);
    consoleDisplay->textsize(12);
    
    consoleInput = new Fl_Input(x + 10, y + h - 30, w - 20, 25);
    consoleInput->callback(consoleInputCallback);
    consoleInput->when(FL_WHEN_ENTER_KEY);
    
    // Приветственное сообщение
    consoleBuffer->text(
        "Добро пожаловать в консоль MixailOS!\n"
        "Введите 'help' для получения списка доступных команд.\n\n"
    );
    
    consoleTab->end();
    return consoleTab;
}

// Создание вкладки с настройками
Fl_Group* createSettingsTab(int x, int y, int w, int h) {
    Fl_Group* settingsTab = new Fl_Group(x, y, w, h, "Настройки");
    
    // Поле для изменения имени пользователя
    new Fl_Box(x + 10, y + 20, 120, 25, "Имя пользователя:");
    usernameInput = new Fl_Input(x + 140, y + 20, 200, 25);
    usernameInput->callback(usernameCallback);
    usernameInput->when(FL_WHEN_ENTER_KEY);
    
    // Загружаем текущее имя пользователя
    char* username = GoGetUsername();
    if (username) {
        usernameInput->value(username);
        freeString(username);
    }
    
    // Кнопка для изменения обоев
    Fl_Button* wallpaperBtn = new Fl_Button(x + 10, y + 60, 150, 25, "Изменить обои");
    wallpaperBtn->callback(changeWallpaperCallback);
    
    settingsTab->end();
    return settingsTab;
}

// Создание вкладки с файловым менеджером
Fl_Group* createFileManagerTab(int x, int y, int w, int h) {
    Fl_Group* fileManagerTab = new Fl_Group(x, y, w, h, "Файлы");
    
    // Текстовое поле с отображением файлов
    Fl_Text_Buffer* fileBuffer = new Fl_Text_Buffer();
    Fl_Text_Display* fileDisplay = new Fl_Text_Display(x + 10, y + 10, w - 20, h - 20);
    fileDisplay->buffer(fileBuffer);
    
    // Загружаем список файлов
    char* fileList = GoGetFileList();
    if (fileList) {
        // Заменяем разделители на новую строку
        std::string files(fileList);
        size_t pos = 0;
        while ((pos = files.find("|", pos)) != std::string::npos) {
            files.replace(pos, 1, "\n");
            pos += 1;
        }
        
        fileBuffer->text(files.c_str());
        freeString(fileList);
    }
    
    fileManagerTab->end();
    return fileManagerTab;
}

// Создание вкладки с браузером
Fl_Group* createBrowserTab(int x, int y, int w, int h) {
    Fl_Group* browserTab = new Fl_Group(x, y, w, h, "Браузер");
    
    // Поле для ввода URL
    browserUrlInput = new Fl_Input(x + 50, y + 10, w - 120, 25, "URL:");
    
    // Кнопка перехода
    Fl_Button* goButton = new Fl_Button(x + w - 60, y + 10, 50, 25, "Go");
    goButton->callback(browserGoCallback);
    
    // Поле для отображения содержимого
    browserBuffer = new Fl_Text_Buffer();
    browserDisplay = new Fl_Text_Display(x + 10, y + 45, w - 20, h - 55);
    browserDisplay->buffer(browserBuffer);
    
    // Приветственное сообщение
    browserBuffer->text(
        "Добро пожаловать в браузер MixailOS!\n"
        "Введите URL в поле выше и нажмите 'Go' для начала работы.\n"
    );
    
    browserTab->end();
    return browserTab;
}

// Создание вкладки с калькулятором
Fl_Group* createCalculatorTab(int x, int y, int w, int h) {
    Fl_Group* calcTab = new Fl_Group(x, y, w, h, "Калькулятор");
    
    // Поле для ввода выражения
    calculatorInput = new Fl_Input(x + 90, y + 10, w - 180, 25, "Выражение:");
    
    // Кнопка вычисления
    Fl_Button* calcButton = new Fl_Button(x + w - 80, y + 10, 70, 25, "Вычислить");
    calcButton->callback(calculateCallback);
    
    // Поле для отображения результата
    calculatorBuffer = new Fl_Text_Buffer();
    calculatorDisplay = new Fl_Text_Display(x + 10, y + 45, w - 20, h - 55);
    calculatorDisplay->buffer(calculatorBuffer);
    
    // Инструкции
    calculatorBuffer->text(
        "Простой калькулятор MixailOS\n"
        "Поддерживаемые операции: +, -, *, /\n"
        "Пример: 2 + 3 * 4\n\n"
    );
    
    calcTab->end();
    return calcTab;
}

// Главная функция для создания интерфейса
extern "C" void RunUI() {
    // Инициализация FLTK
    Fl::scheme("gtk+");
    
    // Инициализация поддержки изображений в FLTK
    fl_register_images();
    
    const int windowWidth = 800;
    const int windowHeight = 600;
    
    // Создание главного окна
    mainWindow = new Fl_Double_Window(windowWidth, windowHeight, "MixailOS");
    
    // Загружаем имя пользователя для заголовка
    char* username = GoGetUsername();
    if (username) {
        std::string title = "MixailOS - " + std::string(username);
        mainWindow->label(title.c_str());
        freeString(username);
    }
    
    // Создаем фон (для обоев)
    wallpaperBox = new Fl_Box(0, 0, windowWidth, windowHeight);
    wallpaperBox->box(FL_FLAT_BOX);
    wallpaperBox->color(FL_WHITE);
    
    // Загружаем обои
    updateWallpaper();
    
    // Создаем вкладки для различных функций
    Fl_Tabs* tabs = new Fl_Tabs(20, 40, windowWidth - 40, windowHeight - 60);
    
    // Добавляем вкладки с различными функциями
    createConsoleTab(20, 70, windowWidth - 40, windowHeight - 90);
    createFileManagerTab(20, 70, windowWidth - 40, windowHeight - 90);
    createBrowserTab(20, 70, windowWidth - 40, windowHeight - 90);
    createCalculatorTab(20, 70, windowWidth - 40, windowHeight - 90);
    createSettingsTab(20, 70, windowWidth - 40, windowHeight - 90);
    
    tabs->end();
    
    mainWindow->end();
    mainWindow->show();
    
    // Запуск главного цикла FLTK
    Fl::run();
} 