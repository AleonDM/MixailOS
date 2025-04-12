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

REM Set Go environment
echo Setting up Go environment...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

REM Download Go dependencies
echo Downloading Go dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo Error downloading Go dependencies
    goto :error
)

echo Tidying up Go modules...
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo Error tidying Go modules
    goto :error
)

REM Creating build directory
if not exist build mkdir build
cd build

REM Set explicit paths to executables
set PATH=C:\msys64\ucrt64\bin;%PATH%

REM Running CMake with MSYS2 MinGW generator
echo Configuring project with CMake...
cmake .. -G "MinGW Makefiles" -DCMAKE_C_COMPILER=C:/msys64/ucrt64/bin/gcc.exe -DCMAKE_CXX_COMPILER=C:/msys64/ucrt64/bin/g++.exe -DCMAKE_MAKE_PROGRAM=C:/msys64/ucrt64/bin/mingw32-make.exe
if %ERRORLEVEL% NEQ 0 (
    echo Error during CMake configuration - trying alternative approach...
    
    REM Try with NMake if available
    where nmake >nul 2>nul
    if %ERRORLEVEL% EQU 0 (
        echo Trying to configure with NMake...
        cmake .. -G "NMake Makefiles" -DCMAKE_C_COMPILER=C:/msys64/ucrt64/bin/gcc.exe -DCMAKE_CXX_COMPILER=C:/msys64/ucrt64/bin/g++.exe
    ) else (
        echo NMake not found. Please run this from a VS Developer Command Prompt.
        cd ..
        goto :error
    )
    
    if %ERRORLEVEL% NEQ 0 (
        cd ..
        goto :error
    )
)

REM Building the project
echo Building project...
IF EXIST mingw32-make.exe (
    mingw32-make
) ELSE IF EXIST C:\msys64\ucrt64\bin\mingw32-make.exe (
    C:\msys64\ucrt64\bin\mingw32-make
) ELSE IF EXIST nmake.exe (
    nmake
) ELSE (
    echo Error: Could not find make tool
    cd ..
    goto :error
)

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
echo 2. Run: pacman -S mingw-w64-ucrt-x86_64-gcc mingw-w64-ucrt-x86_64-make mingw-w64-ucrt-x86_64-cmake mingw-w64-ucrt-x86_64-fltk
echo.
echo For Go dependency issues:
echo 1. Run: go mod download
echo 2. Run: go mod tidy
echo.
echo Make sure C:\msys64\ucrt64\bin is in your PATH environment variable.

:end
echo Press any key to exit...
pause > nul 