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

# Компиляция Go в статическую библиотеку перед CMake
message(STATUS "Compiling Go code into static library...")
execute_process(
    COMMAND go build -buildmode=c-archive -o ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a ${CMAKE_CURRENT_SOURCE_DIR}/main.go
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    RESULT_VARIABLE GO_BUILD_RESULT
)

if(NOT GO_BUILD_RESULT EQUAL 0)
    message(FATAL_ERROR "Failed to compile Go code into static library")
endif()

# Создаем библиотеку из скомпилированного файла .a
add_library(mixailos_lib STATIC IMPORTED)
set_target_properties(mixailos_lib PROPERTIES IMPORTED_LOCATION ${CMAKE_CURRENT_BINARY_DIR}/libmixailos.a)

# Создание исполняемого файла
add_executable(MixailOS ${GUI_TYPE} ${CPP_SOURCES})

# Определяем платформо-зависимые библиотеки
if(UNIX)
    set(PLATFORM_LIBS dl pthread X11)
elseif(WIN32)
    set(PLATFORM_LIBS ws2_32 ole32 uuid)
endif()

# Связывание библиотек
target_link_libraries(MixailOS PRIVATE mixailos_lib ${FLTK_LIBRARIES} ${PLATFORM_LIBS})

# Установка
install(TARGETS MixailOS DESTINATION bin)

# Дополнительная информация при компиляции
message(STATUS "FLTK include directories: ${FLTK_INCLUDE_DIR}")
message(STATUS "FLTK libraries: ${FLTK_LIBRARIES}") 