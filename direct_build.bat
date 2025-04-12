@echo on
echo Building MixailOS with direct compilation...

REM Setting environment variables
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
set CC=gcc
set CXX=g++

REM Clean previous build
if exist build rmdir /S /Q build
mkdir build

REM First, compile Go code to create a shared library
echo Compiling Go code to shared library...
go build -v -x -buildmode=c-shared -o build\mixailos.dll main.go

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Go code
    goto :error
)

REM Check DLL exports
echo DLL exports:
nm -g build\mixailos.dll | findstr Go

REM Create a simple C++ app to test the DLL
echo Creating a test file...
cd build
echo #include ^<stdio.h^> > test.cpp
echo extern "C" { >> test.cpp
echo     void InitMixailOS(); >> test.cpp
echo     void StartUI(); >> test.cpp
echo     char* GetConfigUsername(); >> test.cpp
echo } >> test.cpp
echo int main() { >> test.cpp
echo     printf("Testing DLL functions...\n"); >> test.cpp
echo     InitMixailOS(); >> test.cpp
echo     char* username = GetConfigUsername(); >> test.cpp
echo     printf("Username: %s\n", username); >> test.cpp
echo     // StartUI() would launch the UI, so we don't call it here >> test.cpp
echo     return 0; >> test.cpp
echo } >> test.cpp

REM Compile and link the test application
g++ -o test.exe test.cpp -L. -lmixailos

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile test application
    cd ..
    goto :error
)

REM Create the WinMain entry point
echo Creating WinMain.cpp...
echo #include ^<windows.h^> > winmain.cpp
echo extern "C" { >> winmain.cpp
echo     void InitMixailOS(); >> winmain.cpp
echo     void StartUI(); >> winmain.cpp
echo } >> winmain.cpp
echo int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) { >> winmain.cpp
echo     InitMixailOS(); >> winmain.cpp
echo     StartUI(); >> winmain.cpp
echo     return 0; >> winmain.cpp
echo } >> winmain.cpp
echo int main(int argc, char** argv) { >> winmain.cpp
echo     return WinMain(GetModuleHandle(NULL), NULL, GetCommandLine(), SW_SHOW); >> winmain.cpp
echo } >> winmain.cpp

REM Compile UI source
echo Compiling ui.cpp...
g++ -c -I.. -I"C:/msys64/ucrt64/include" -I"C:/msys64/ucrt64/include/FL" ..\ui\ui.cpp -o ui.o

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile ui.cpp
    cd ..
    goto :error
)

REM Compile main application
echo Compiling main application...
g++ -o MixailOS.exe winmain.cpp ui.o -I.. -I"C:/msys64/ucrt64/include" -I"C:/msys64/ucrt64/include/FL" -L. -lmixailos -L"C:/msys64/ucrt64/lib" -lfltk -lfltk_images -lgdi32 -lole32 -luuid -lcomdlg32 -lws2_32 -mwindows

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile main application
    cd ..
    goto :error
)

echo Testing shared library...
.\test.exe

echo Build successful! You can now run build\MixailOS.exe
cd ..
goto :end

:error
echo Build failed. Please check the error messages above.

:end
echo Press any key to exit...
pause > nul 