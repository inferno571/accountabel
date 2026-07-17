#!/bin/bash

echo "Building Accountabel AI for Linux/macOS (Production)..."

# -s omits the symbol table and debug information
# -w omits the DWARF symbol table
go build -ldflags "-s -w" -o AccountabelAI main.go

if [ $? -eq 0 ]; then
    echo "Build successful! Executable created as AccountabelAI"
else
    echo "Build failed."
fi
