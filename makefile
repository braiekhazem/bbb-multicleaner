run: 
	go run *.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bbb-multicleaner main.go logging.go

build:
	go build

clean:
	rm -f bbb-multicleaner