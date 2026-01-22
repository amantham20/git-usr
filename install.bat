@echo off
REM Installation script for git-usr (Windows)

echo Installing git-usr...

REM Get the directory where this script is located
set SCRIPT_DIR=%~dp0

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Go is not installed. Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Build the binary
echo Building git-usr...
cd /d "%SCRIPT_DIR%"
go build -o git-usr.exe main.go

if %ERRORLEVEL% NEQ 0 (
    echo Build failed
    pause
    exit /b 1
)

REM Add to PATH (user environment variable)
setx PATH "%PATH%;%SCRIPT_DIR%"

echo.
echo Successfully installed git-me!
echo.
echo Usage (as a git subcommand):
echo   git usr work              # Switch to work profile
echo   git usr personal          # Switch to personal profile
echo   git usr list              # List all profiles
echo   git usr add ^<profile^>     # Add a new profile
echo.
echo First time setup:
echo   1. Restart your terminal/command prompt
echo   2. Run 'git usr list' to see default profiles
echo   3. Update profiles with: git usr add work "Your Name" "email@example.com"
echo   4. Switch profiles with: git usr work
echo.
echo Shell completion (optional):
echo   git usr completion powershell ^> git-usr-completion.ps1
echo   Then add to your $PROFILE: . path\to\git-usr-completion.ps1
echo.
echo Config stored in: %%APPDATA%%\git-usr\profiles.json
