@ECHO OFF

set exec_name=suppyship

echo Compiling executable ...
go build -o "%exec_name%.exe" -ldflags="-H=windowsgui -s -w" .
echo.

:: Fail if compilation was unsuccesful
if %ERRORLEVEL% NEQ 0 (
    exit
)

:: Compress executable
echo Compressing executable ...
"bin/upx" -v -9 -o "%exec_name%-c.exe" %exec_name%.exe
del "%exec_name%.exe">nul
move "%exec_name%-c.exe" "%exec_name%.exe" >nul