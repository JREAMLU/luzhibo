@echo off
title Auto Make

set FNAME=luzhibo
set PNAME=%FNAME%
set GPATH=https://github.com/Baozisoftware/Luzhibo-go.git
set CPATH=%cd%


::init
echo Initing...
go get github.com\PuerkitoBio\goquery
go get github.com\pkg\browser
go get github.com\Baozisoftware\qrcode-terminal-go
if "%1%"=="init" goto done

if "%GOPATH%"=="" set GOPATH=%UserProfile%\go
set SPATH=%GOPATH%\src\%PNAME%
git clone %GPATH% %SPATH%
cd %SPATH%
git pull
cd %CPATH%

if exist releases rd /s /q releases
md releases

::386:7
set GOARCH=386

set GOOS=darwin
call:make
set GOOS=freebsd
call:make
set GOOS=linux
call:make
set GOOS=netbsd
call:make
set GOOS=openbsd
call:make
set GOOS=plan9
call:make
set GOOS=windows
call:make

::amd64:9
set GOARCH=amd64

set GOOS=darwin
call:make
set GOOS=dragonfly
call:make
set GOOS=freebsd
call:make
set GOOS=linux
call:make
set GOOS=netbsd
call:make
set GOOS=openbsd
call:make
set GOOS=plan9
call:make
set GOOS=solaris
call:make
set GOOS=windows
call:make

::arm:6
set GOARCH=arm

set GOOS=android
call:make
set GOOS=darwin
call:make
set GOOS=freebsd
call:make
set GOOS=linux
call:make
set GOOS=netbsd
call:make
set GOOS=openbsd
call:make

::arm64:2
set GOARCH=arm64

set GOOS=darwin
call:make
set GOOS=linux
call:make

::mips:1
set GOARCH=mips

set GOOS=linux
call:make

::mipsle:1
set GOARCH=mipsle

set GOOS=linux
call:make

::mips64:1
set GOARCH=mips64

set GOOS=linux
call:make

::mips64le:1
set GOARCH=mips64le

set GOOS=linux
call:make

::ppc64:1
set GOARCH=ppc64

set GOOS=linux
call:make

::ppc64le:1
set GOARCH=ppc64le

set GOOS=linux
call:make

cd releases
..\7z a -t7z ..\releases.7z -r -mx=9 -m0=LZMA2 -ms=100m -mf=on -mhc=on -mmt=on

:done
echo All done.
pause
goto:eof

:make
set TNAME=%FNAME%_%GOOS%_%GOARCH%
if %GOOS%==windows set TNAME=%TNAME%.exe
set TPATH=releases\%TNAME%
echo Building %TNAME%....
go build -ldflags="-s -w" -o %TPATH% %PNAME%
upx --best -q %TPATH%
goto:eof
