BINARY_NAME=bin/chatgpt-cli

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o ${BINARY_NAME}-amd64.exe main.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-x64.exe
