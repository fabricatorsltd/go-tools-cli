@echo off
cd /d "C:\Users\PietrodiCaprio\Dev\repos\go-tools\go-tools-cli"
echo Running go mod tidy...
go mod tidy
echo.
echo TIDY_RC=%ERRORLEVEL%
echo.
echo Running go build ./...
go build ./...
echo.
echo BUILD_RC=%ERRORLEVEL%
