#!/bin/bash

echo "Building MixailOS with pure Go..."

# Проверка Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go to continue."
    exit 1
fi

go version

# Загрузка зависимостей
echo "Installing dependencies..."
go mod tidy

if [ $? -ne 0 ]; then
    echo "Error downloading dependencies"
    exit 1
fi

# Компиляция проекта
echo "Building MixailOS..."
go build -o MixailOS

if [ $? -ne 0 ]; then
    echo "Error building project"
    exit 1
fi

echo "MixailOS successfully built!"
echo "Run ./MixailOS to start the application." 