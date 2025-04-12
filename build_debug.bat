@echo on
echo Building MixailOS in Debug Mode...

REM Checking Go installation
where go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go to continue.
    goto :error
)
go version

REM Checking CMake installation
where cmake
if %ERRORLEVEL% NEQ 0 (
    echo Error: CMake is not installed. Please install CMake to continue.
    goto :error
)
cmake --version

REM Checking MinGW installation
where gcc
if %ERRORLEVEL% NEQ 0 (
    echo Error: MinGW gcc is not installed or not in PATH.
    echo Please install MSYS2 and add C:\msys64\ucrt64\bin to your PATH.
    goto :error
)
gcc --version

REM Set Go environment with verbose output
echo Setting up Go environment...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
echo CGO_ENABLED=%CGO_ENABLED%
echo GOOS=%GOOS%
echo GOARCH=%GOARCH%

REM Show go.mod content
echo Go.mod content:
type go.mod

REM Clean previous build artifacts
echo Cleaning previous build...
if exist build\MixailOS.exe del /F build\MixailOS.exe
if exist build\libmixailos.* del /F build\libmixailos.*

REM Test Go compilation separately first
echo Testing Go compilation...
go build -v -x -buildmode=c-archive -o test_lib.a main.go
if %ERRORLEVEL% NEQ 0 (
    echo Error during Go test compilation
    goto :error
)
echo Go test compilation successful!
del /F test_lib.a test_lib.h

REM Creating build directory
if not exist build mkdir build
cd build

REM Running CMake with MSYS2 MinGW generator with verbose output
echo Configuring project with CMake...
cmake .. -G "MinGW Makefiles" -DCMAKE_C_COMPILER=C:/msys64/ucrt64/bin/gcc.exe -DCMAKE_CXX_COMPILER=C:/msys64/ucrt64/bin/g++.exe -DCMAKE_MAKE_PROGRAM=C:/msys64/ucrt64/bin/mingw32-make.exe -DCMAKE_BUILD_TYPE=Debug -DCMAKE_VERBOSE_MAKEFILE=ON
if %ERRORLEVEL% NEQ 0 (
    echo Error during CMake configuration
    cd ..
    goto :error
)

REM Building the project
echo Building project...
mingw32-make VERBOSE=1
if %ERRORLEVEL% NEQ 0 (
    echo Error during build
    cd ..
    goto :error
)

REM Check if the binary was created
if not exist MixailOS.exe (
    echo Error: Could not find MixailOS.exe in the build directory
    cd ..
    goto :error
)

echo Checking dependencies...
objdump -p MixailOS.exe

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
echo For Go dependency issues, run:
echo   go get -u github.com/AleonDM/MixailOS/...
echo   go mod download
echo   go mod tidy
echo.
echo Make sure C:\msys64\ucrt64\bin is in your PATH environment variable.

:end
echo Press any key to exit...
pause > nul 