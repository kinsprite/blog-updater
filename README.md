# Blog Updater in Golang

## Build

on Windows:
```cmd
set GOPROXY=https://goproxy.io
go mod vendor
go build -mod=vendor -o blog-updater.exe .
```

or Linux:

```shell
export GOPROXY=https://goproxy.io
go mod vendor
go build -mod=vendor -o blog-updater .
```

## Running Env

```shell
export GIN_MODE=release
export LISTENING_ADDRESS=127.0.0.1:8080
export SERVER_SIGNATURE=sha1=xxx-yyy-zzz
export SHELL_SCRIPT_FILE=/etc/blog-updater/do-update.sh
./blog-updater
```
