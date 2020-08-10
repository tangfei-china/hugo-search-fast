# Mac上编译，Linux上执行
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o search_index
