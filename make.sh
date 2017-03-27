#!/bin/sh

echo -n "\033]0;Auto Make\007"

FNAME=luzhibo
PNAME=github.com/Baozisoftware/$FNAME
GPATH=https://github.com/Baozisoftware/luzhibo.git
CPATH=`pwd`
BPATH=`dirname $0`
UPX=$BPATH/upx
chmod +x $UPX

MAKE()
{
	TNAME="$FNAME"_"$GOOS"_"$GOARCH"
	if [ "$GOOS" = "windows" ]; then
		TNAME=$TNAME.exe
       		GOOS=$GOOS GOARCH=$GOARCH go generate $PNAME
		GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $TPATH $PNAME
	fi
	TPATH=releases/$TNAME
	echo Building $TNAME....
	#GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $TPATH $PNAME
    if [ -f "$SPATH/resource.syso" ]; then
        rm $SPATH/resource.syso
    fi
	#$BPATH/upx --best -q $TPATH
	
}

DONE()
{
	echo All done.
	exit 0
}

#init
echo Initing...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get github.com/PuerkitoBio/goquery
go get github.com/pkg/browser
go get github.com/Baozisoftware/qrcode-terminal-go
go get github.com/mattn/go-isatty
if [ "$1" = "init" ]; then
	DONE
fi
if [ "$GOPATH" = "" ]; then 
	GOPATH=~/go
fi
PATH=$PATH:$GOPATH/bin
SPATH=$GOPATH/src/$PNAME
git clone $GPATH $SPATH
cd $SPATH
git pull
cd $CPATH


if [ -d releases ]; then
	rm -rf releases
fi
mkdir releases
#386:7
GOARCH=386

GOOS=darwin
MAKE
GOOS=freebsd
MAKE
GOOS=linux
MAKE
GOOS=netbsd
MAKE
GOOS=openbsd
MAKE
GOOS=plan9
MAKE
GOOS=windows
MAKE

#amd64:9
GOARCH=amd64

GOOS=darwin
MAKE
GOOS=dragonfly
MAKE
GOOS=freebsd
MAKE
GOOS=linux
MAKE
GOOS=netbsd
MAKE
GOOS=openbsd
MAKE
GOOS=plan9
MAKE
GOOS=solaris
MAKE
GOOS=windows
MAKE

#arm:6
GOARCH=arm

GOOS=android
MAKE
GOOS=darwin
MAKE
GOOS=freebsd
MAKE
GOOS=linux
MAKE
GOOS=netbsd
MAKE
GOOS=openbsd
MAKE

#arm64:2
GOARCH=arm64

GOOS=darwin
MAKE
GOOS=linux
MAKE

#mips:1
GOARCH=mips

GOOS=linux
MAKE

#mipsle:1
GOARCH=mipsle

GOOS=linux
MAKE

#mips64:1
GOARCH=mips64

GOOS=linux
MAKE

#mips64le:1
GOARCH=mips64le

GOOS=linux
MAKE

#ppc64:1
GOARCH=ppc64

GOOS=linux
MAKE

#ppc64le:1
GOARCH=ppc64le

GOOS=linux
MAKE

DONE

