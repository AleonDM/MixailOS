@echo off
echo Building MixailOS...

REM Checking Go installation
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go to continue.
    goto :error
)

REM Checking CMake installation
where cmake >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: CMake is not installed. Please install CMake to continue.
    goto :error
)

REM Checking MinGW installation
where gcc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: MinGW gcc is not installed or not in PATH.
    echo Please install MSYS2 and add C:\msys64\mingw64\bin to your PATH.
    goto :error
)

REM Creating build directory
if not exist build mkdir build
cd build

REM Running CMake
echo Configuring project with CMake...
cmake .. -G "MinGW Makefiles" -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ -DCMAKE_MAKE_PROGRAM=mingw32-make
if %ERRORLEVEL% NEQ 0 (
    echo Error during CMake configuration
    cd ..
    goto :error
)

REM Building the project
echo Building project...
mingw32-make
if %ERRORLEVEL% NEQ 0 (
    echo Error during build
    cd ..
    goto :error
)

echo Build completed successfully!
echo Run .\MixailOS.exe to start MixailOS
cd ..
goto :end

:error
echo.
echo Build failed. See the error messages above.
echo.
echo You might need to install MSYS2 and then run:
echo pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake mingw-w64-x86_64-make mingw-w64-x86_64-fltk
echo.
echo And make sure to add C:\msys64\mingw64\bin to your PATH.

:end
echo Press any key to exit...
pause > nul 