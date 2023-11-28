dev:
	go mod tidy
	go build -o bin/cmd cmd/main.go
	bin/cmd