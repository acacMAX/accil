@echo off
go build -o accil.exe .
if not exist "%USERPROFILE%\.accil\bin" mkdir "%USERPROFILE%\.accil\bin"
copy /y accil.exe "%USERPROFILE%\.accil\bin\" >nul
echo Build and install complete!
echo Installed to: %USERPROFILE%\.accil\bin\accil.exe
