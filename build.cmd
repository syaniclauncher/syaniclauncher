@echo off
setlocal

if not exist version.txt (
    echo version.txt not found
    exit /b 1
)

set /p VERSION=<version.txt
if not defined VERSION (
    echo version.txt is empty
    exit /b 1
)

for /f %%i in ('git rev-parse --short HEAD') do set "COMMIT=%%i"
if not defined COMMIT (
    echo Could not determine the current git commit
    exit /b 1
)

set "BUILD_ID=%VERSION%:%COMMIT%"
wails build -ldflags "-X main.BUILD_ID=%BUILD_ID%"
exit /b %ERRORLEVEL%
