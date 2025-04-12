@echo on
echo Building MixailOS in Debug Mode...

REM Проверка установки Go
where go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go to continue.
    goto :error
)
go version

REM Проверка установки CMake
where cmake
if %ERRORLEVEL% NEQ 0 (
    echo Error: CMake is not installed. Please install CMake to continue.
    goto :error
)
cmake --version

REM Проверка установки MinGW
where gcc
if %ERRORLEVEL% NEQ 0 (
    echo Error: MinGW gcc is not installed or not in PATH.
    echo Please install MSYS2 and add C:\msys64\ucrt64\bin to your PATH.
    goto :error
)
gcc --version

REM Настройка переменных окружения Go
echo Setting up Go environment...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
echo CGO_ENABLED=%CGO_ENABLED%
echo GOOS=%GOOS%
echo GOARCH=%GOARCH%

REM Проверка содержимого go.mod
echo Go.mod content:
type go.mod

REM Очистка предыдущих артефактов сборки
echo Cleaning previous build...
if exist build rmdir /S /Q build
mkdir build
cd build

REM Компиляция Go кода в статическую библиотеку с подробным выводом
echo Compiling Go code to static library...
go build -v -x -buildmode=c-archive -o libmixailos.a ..\main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error during Go compilation
    cd ..
    goto :error
)

REM Проверка наличия скомпилированной библиотеки
if not exist libmixailos.a (
    echo Error: Failed to create libmixailos.a
    cd ..
    goto :error
)

REM Попробуем вручную запустить nm для просмотра символов
echo Checking exported symbols in libmixailos.a...
nm -g libmixailos.a | findstr Go

REM Запуск CMake с генератором MinGW Makefiles
echo Configuring project with CMake...
cmake .. -G "MinGW Makefiles" -DCMAKE_C_COMPILER=C:/msys64/ucrt64/bin/gcc.exe -DCMAKE_CXX_COMPILER=C:/msys64/ucrt64/bin/g++.exe -DCMAKE_MAKE_PROGRAM=C:/msys64/ucrt64/bin/mingw32-make.exe -DCMAKE_BUILD_TYPE=Debug -DCMAKE_VERBOSE_MAKEFILE=ON
if %ERRORLEVEL% NEQ 0 (
    echo Error during CMake configuration
    cd ..
    goto :error
)

REM Компиляция проекта
echo Building project...
mingw32-make VERBOSE=1
if %ERRORLEVEL% NEQ 0 (
    echo Error during build
    cd ..
    goto :error
)

REM Проверка, был ли создан исполняемый файл
if not exist MixailOS.exe (
    echo Error: Failed to create MixailOS.exe
    cd ..
    goto :error
)

REM Проверка зависимостей исполняемого файла
echo Checking dependencies...
objdump -p MixailOS.exe

echo Build completed successfully!
echo Run .\build\MixailOS.exe to start MixailOS
cd ..
goto :end

:error
echo.
echo Build failed. See the error messages above.
echo.
echo Проверка на общие проблемы:
echo - Убедитесь, что CGO_ENABLED=1
echo - Попробуйте очистить кэш Go модулей: go clean -cache -modcache
echo - Возможно, не хватает библиотек FLTK или других зависимостей
echo.
echo Обновление пакетов MSYS2:
echo 1. Откройте терминал MSYS2 UCRT64
echo 2. Выполните: pacman -Syu
echo 3. Выполните: pacman -S mingw-w64-ucrt-x86_64-gcc mingw-w64-ucrt-x86_64-make mingw-w64-ucrt-x86_64-cmake mingw-w64-ucrt-x86_64-fltk

:end
echo Press any key to exit...
pause > nul 