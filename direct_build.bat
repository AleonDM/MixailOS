@echo on
echo Building MixailOS with direct compilation...

REM Setting environment variables
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

REM Ensure the build directory exists
if not exist build mkdir build

REM First, compile Go code to create a library
echo Compiling Go code...
go build -buildmode=c-shared -o build/mixailos.dll main.go

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Go code
    goto :error
)

REM Check DLL exports
echo DLL exports:
dumpbin /exports build/mixailos.dll > build/exports.txt
type build\exports.txt

REM Now compile the C++ code
echo Compiling C++ code...
cd build

REM Create the WinMain entry point
echo Creating WinMain.cpp...
echo #include ^<windows.h^> > winmain.cpp
echo #include ^<FL/Fl.H^> >> winmain.cpp
echo extern "C" { >> winmain.cpp
echo     void StartUI(); >> winmain.cpp
echo } >> winmain.cpp
echo int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) { >> winmain.cpp
echo     StartUI(); >> winmain.cpp
echo     return 0; >> winmain.cpp
echo } >> winmain.cpp
echo int main(int argc, char** argv) { >> winmain.cpp
echo     return WinMain(GetModuleHandle(NULL), NULL, GetCommandLine(), SW_SHOW); >> winmain.cpp
echo } >> winmain.cpp

REM Compile C++ code with direct g++ command
g++ -o MixailOS.exe winmain.cpp ../ui/ui.cpp -I.. -I"C:/msys64/ucrt64/include" -I"C:/msys64/ucrt64/include/FL" -L. -lmixailos -L"C:/msys64/ucrt64/lib" -lfltk -lfltk_images -lgdi32 -lole32 -luuid -lcomctl32 -lcomdlg32 -lws2_32 -mwindows

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile C++ code
    cd ..
    goto :error
)

echo Build successful! You can now run build\MixailOS.exe
cd ..
goto :end

:error
echo Build failed. Please check the error messages above.

:end
echo Press any key to exit...
pause > nul 