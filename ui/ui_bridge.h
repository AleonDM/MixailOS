#ifndef UI_BRIDGE_H
#define UI_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

// Функции для модуля CGO
void FreeString(char* s);

// Функции экспортированные из main.go
void InitMixailOS();
void StartUI();
char* GetConfigUsername();

// Функции экспортированные из ui/cgo_bridge.go
char* GoGetUsername();
void GoSetUsername(char* username);
char* GoGetCurrentDirectory();
char* GoExecuteConsoleCommand(char* cmd);
char* GoGetFileList();
void GoChangeWallpaper(char* path);
char* GoGetWallpaperPath();
char* GoCreateTextFile(char* name, char* content);
char* GoReadTextFile(char* name);
void RunUI();

// Функция C++ для запуска UI
void RunUI();

#ifdef __cplusplus
}
#endif

#endif // UI_BRIDGE_H 