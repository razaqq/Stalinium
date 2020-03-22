# Stalinium Installer

## Build steps

#### 1. Get qt binding
1. set GO111MODULE=off
2. go get -v github.com/therecipe/qt/cmd/...
3. %GOPATH%\bin\qtsetup test
4. %GOPATH%\bin\qtsetup -test=false
5. %GOPATH%\bin\qtenv.bat
6. setx PATH "%PATH%"

#### 2. Generate static data
7. qtmoc
8. qtrcc
9. go get github.com/akavel/rsrc
10. rsrc -arch amd64 -manifest StaliniumInstaller.exe.manifest -ico icon.ico -o StaliniumInstaller.syso

#### 3. Build
11. qtdeploy -docker build windows_64_static

## GOENV (go env)
```console
set GO111MODULE=off
set GOARCH=amd64
set GOBIN=
set GOCACHE=C:\Users\<User>\AppData\Local\go-build
set GOENV=C:\Users\<User>\AppData\Roaming\go\env
set GOEXE=.exe
set GOFLAGS=
set GOHOSTARCH=amd64
set GOHOSTOS=windows
set GONOPROXY=
set GONOSUMDB=
set GOOS=windows
set GOPATH=C:\Users\<User>\go
set GOPRIVATE=
set GOPROXY=https://proxy.golang.org,direct
set GOROOT=C:\Program Files (x86)\Go
set GOSUMDB=sum.golang.org
set GOTMPDIR=
set GOTOOLDIR=C:\Program Files (x86)\Go\pkg\tool\windows_amd64
set GCCGO=gccgo
set AR=ar
set CC=gcc
set CXX=g++
set CGO_ENABLED=1
set GOMOD=
set CGO_CFLAGS=-g -O2
set CGO_CPPFLAGS=
set CGO_CXXFLAGS=-g -O2
set CGO_FFLAGS=-g -O2
set CGO_LDFLAGS=-g -O2
set PKG_CONFIG=pkg-config
set GOGCCFLAGS=-m64 -mthreads -fmessage-length=0 -fdebug-prefix-map=C:\Users\<User>\AppData\Local\Temp\go-build480071318=/tmp/go-build -gno-record-gcc-switches
````