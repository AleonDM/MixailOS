#!/bin/bash

echo "Сборка MixailOS..."

# Проверка Go
if ! command -v go &> /dev/null; then
    echo "Ошибка: Go не установлен. Пожалуйста, установите Go для продолжения."
    exit 1
fi

# Проверка CMake
if ! command -v cmake &> /dev/null; then
    echo "Ошибка: CMake не установлен. Пожалуйста, установите CMake для продолжения."
    exit 1
fi

# Создание и переход в директорию сборки
if [ ! -d "build" ]; then
    mkdir build
fi

cd build

# Запуск CMake
echo "Настройка проекта с помощью CMake..."
cmake .. || { echo "Ошибка при конфигурации CMake"; exit 1; }

# Запуск сборки
echo "Сборка проекта..."
make || { echo "Ошибка при сборке"; exit 1; }

echo "Сборка успешно завершена!"
echo "Запустите ./build/MixailOS для запуска MixailOS" 