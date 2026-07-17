@echo off
echo Building Accountabel AI for Windows (Production)...

:: -s omits the symbol table and debug information
:: -w omits the DWARF symbol table
go build -ldflags "-s -w" -o AccountabelAI.exe main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Executable created as AccountabelAI.exe
) else (
    echo Build failed.
)
pause
