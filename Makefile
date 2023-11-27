dev:
	go mod tidy
	go build -o bin/bin cmd/main.go
	bin/bin