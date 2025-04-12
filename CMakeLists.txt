cmake_minimum_required(VERSION 3.10)
project(MixailOS)

# Установка стандарта C++
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Поиск библиотек
find_package(FLTK REQUIRED)
include_directories(${FLTK_INCLUDE_DIR})

# Включение директорий
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

# Исходные файлы C++
file(GLOB CPP_SOURCES "ui/*.cpp")

# Компиляция Go в статическую библиотеку
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a
    COMMAND go build -buildmode=c-archive -o ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a ${CMAKE_CURRENT_SOURCE_DIR}/main.go
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/main.go
)

# Создание исполняемого файла
add_executable(MixailOS ${CPP_SOURCES} ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a)

# Связывание библиотек
target_link_libraries(MixailOS ${FLTK_LIBRARIES} dl pthread)

# Если компиляция на Unix-подобных системах
if(UNIX)
    target_link_libraries(MixailOS X11)
endif()

# Установка
install(TARGETS MixailOS DESTINATION bin)

# Дополнительная информация при компиляции
message(STATUS "FLTK include directories: ${FLTK_INCLUDE_DIR}")
message(STATUS "FLTK libraries: ${FLTK_LIBRARIES}") 