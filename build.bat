@echo off
setlocal
set "ACCIL_VER=1.4.0"
if exist VERSION for /f "usebackq delims=" %%a in ("VERSION") do set "ACCIL_VER=%%a"
set "ACCIL_VER=%ACCIL_VER: =%"
go build -buildvcs=false -ldflags="-X github.com/accil/accil/cmd.Version=%ACCIL_VER%" -o accil.exe .
if not exist "%USERPROFILE%\.accil\bin" mkdir "%USERPROFILE%\.accil\bin"
copy /y accil.exe "%USERPROFILE%\.accil\bin\" >nul
echo Build and install complete!
echo Installed to: %USERPROFILE%\.accil\bin\accil.exe
