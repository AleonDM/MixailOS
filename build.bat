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
    echo Please install MSYS2 and add C:\msys64\ucrt64\bin to your PATH.
    goto :error
)

REM Creating build directory
if not exist build mkdir build
cd build

REM Running CMake
echo Configuring project with CMake...
cmake .. -G "MSYS Makefiles" -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ -DCMAKE_MAKE_PROGRAM=make
if %ERRORLEVEL% NEQ 0 (
    echo Error during CMake configuration
    cd ..
    goto :error
)

REM Building the project
echo Building project...
make
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
echo You might need to install or update MSYS2 packages:
echo 1. Open MSYS2 UCRT64 terminal
echo 2. Run: pacman -S mingw-w64-ucrt-x86_64-gcc mingw-w64-ucrt-x86_64-make mingw-w64-ucrt-x86_64-fltk make
echo.
echo And make sure to add C:\msys64\ucrt64\bin to your PATH.

:end
echo Press any key to exit...
pause > nul 