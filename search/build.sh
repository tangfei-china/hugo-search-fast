# Mac 下执行,Linux上执行
env CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o search
