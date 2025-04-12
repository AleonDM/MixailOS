@echo on
echo Building MixailOS with pure Go...

REM Проверка установки Go
where go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go to continue.
    goto :error
)
go version

REM Загружаем необходимые зависимости
echo Installing dependencies...
go mod tidy

if %ERRORLEVEL% NEQ 0 (
    echo Error downloading dependencies
    goto :error
)

REM Компиляция проекта
echo Building MixailOS...
go build -o MixailOS.exe

if %ERRORLEVEL% NEQ 0 (
    echo Error building project
    goto :error
)

echo MixailOS successfully built!
echo Run MixailOS.exe to start the application.
goto :end

:error
echo Build failed. See the error messages above.

:end
echo Press any key to exit...
pause > nul 