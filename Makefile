all:
	go mod tidy
	go build
	gofmt -l -s -w .
