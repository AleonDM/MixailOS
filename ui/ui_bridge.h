#ifndef UI_BRIDGE_H
#define UI_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

// Экспортируемые функции из Go
extern char* GoGetUsername();
extern void GoSetUsername(char* username);
extern char* GoGetCurrentDirectory();
extern char* GoExecuteConsoleCommand(char* cmd);
extern char* GoGetFileList();
extern void GoChangeWallpaper(char* path);
extern char* GoGetWallpaperPath();
extern char* GoCreateTextFile(char* name, char* content);
extern char* GoReadTextFile(char* name);

// Функция для запуска UI, реализованная в C++
extern void RunUI();

#ifdef __cplusplus
}
#endif

#endif // UI_BRIDGE_H 