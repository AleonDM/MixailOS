cmake_minimum_required(VERSION 3.10)
project(MixailOS)

# Установка стандарта C++
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Для Windows - нужно определить точку входа
if(WIN32)
    set(GUI_TYPE WIN32)
endif()

# Поиск библиотек
find_package(FLTK REQUIRED)
include_directories(${FLTK_INCLUDE_DIR})

# Включение директорий
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

# Исходные файлы C++
file(GLOB CPP_SOURCES "ui/*.cpp")

# Определяем файл для WinMain
if(WIN32)
    set(WIN_MAIN_SOURCE "${CMAKE_CURRENT_SOURCE_DIR}/ui/winmain.cpp")
    file(WRITE ${WIN_MAIN_SOURCE} "
        #include <windows.h>
        #include <FL/Fl.H>
        
        extern \"C\" void RunUI();
        
        int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) {
            RunUI();
            return 0;
        }
        
        int main(int argc, char** argv) {
            return WinMain(GetModuleHandle(NULL), NULL, GetCommandLine(), SW_SHOW);
        }
    ")
    list(APPEND CPP_SOURCES ${WIN_MAIN_SOURCE})
endif()

# Компиляция Go в статическую библиотеку
add_custom_command(
    OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a
    COMMAND go build -buildmode=c-archive -o ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a ${CMAKE_CURRENT_SOURCE_DIR}/main.go
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/main.go ${CMAKE_CURRENT_SOURCE_DIR}/core/config.go ${CMAKE_CURRENT_SOURCE_DIR}/core/console.go ${CMAKE_CURRENT_SOURCE_DIR}/core/filesystem.go ${CMAKE_CURRENT_SOURCE_DIR}/ui/cgo_bridge.go
)

# Создание исполняемого файла
add_executable(MixailOS ${GUI_TYPE} ${CPP_SOURCES})

# Добавляем зависимость от libmixailos.a
add_dependencies(MixailOS ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a)

# Определяем платформо-зависимые библиотеки
if(UNIX)
    set(PLATFORM_LIBS dl pthread X11)
elseif(WIN32)
    set(PLATFORM_LIBS ws2_32 ole32 uuid)
endif()

# Связывание библиотек
target_link_libraries(MixailOS PRIVATE ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a ${FLTK_LIBRARIES} ${PLATFORM_LIBS})

# Установка
install(TARGETS MixailOS DESTINATION bin)

# Дополнительная информация при компиляции
message(STATUS "FLTK include directories: ${FLTK_INCLUDE_DIR}")
message(STATUS "FLTK libraries: ${FLTK_LIBRARIES}") 