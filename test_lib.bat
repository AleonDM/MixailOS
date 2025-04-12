@echo on
echo Building and testing Go library...

REM Setting environment variables
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

REM Clean test directory
if exist test_lib rmdir /S /Q test_lib
mkdir test_lib
cd test_lib

REM Compile Go code to DLL
echo Compiling Go code to DLL...
go build -v -buildmode=c-shared -o mixailos.dll ..\main.go

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile Go code to DLL
    cd ..
    goto :error
)

REM Create a simple test application in C
echo Creating test application...
echo #include <stdio.h> > test.c
echo extern void InitMixailOS(); >> test.c
echo extern char* GetConfigUsername(); >> test.c
echo extern void FreeString(char* s); >> test.c
echo int main() { >> test.c
echo     printf("Initializing MixailOS...\n"); >> test.c
echo     InitMixailOS(); >> test.c
echo     char* username = GetConfigUsername(); >> test.c
echo     printf("Username: %s\n", username); >> test.c
echo     FreeString(username); >> test.c
echo     printf("Test completed\n"); >> test.c
echo     return 0; >> test.c
echo } >> test.c

REM Compile test application
echo Compiling test application...
gcc -o test.exe test.c -L. -lmixailos

if %ERRORLEVEL% NEQ 0 (
    echo Failed to compile test application
    cd ..
    goto :error
)

REM Run test application
echo Running test application...
.\test.exe

cd ..
goto :end

:error
echo Test failed. See error messages above.

:end
echo Press any key to exit...
pause > nul 