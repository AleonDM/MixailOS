#include <stdio.h>
#include <string.h>
#include "ui_bridge.h"

// Для тестирования мы будем использовать упрощенную версию UI
void RunUI() {
    printf("RunUI called from C++\n");
    
    // Получаем имя пользователя
    char* username = GoGetUsername();
    printf("Username from Go: %s\n", username);
    FreeString(username);
    
    // Получаем текущую директорию
    char* currentDir = GoGetCurrentDirectory();
    printf("Current directory: %s\n", currentDir);
    FreeString(currentDir);
    
    // Получаем и показываем список файлов
    char* fileList = GoGetFileList();
    printf("Files:\n%s\n", fileList);
    FreeString(fileList);
    
    // Выполняем команду консоли
    char* cmdResult = GoExecuteConsoleCommand("help");
    printf("Console command result:\n%s\n", cmdResult);
    FreeString(cmdResult);
    
    printf("UI Test completed\n");
} 