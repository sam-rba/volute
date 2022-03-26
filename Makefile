build: tidy format
	go build
run: tidy format
	go run .

format:
	gofmt -s -w .

tidy:
	go mod tidy
