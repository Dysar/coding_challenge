dev: swag
	go mod tidy
	go build -o bin/bin cmd/main.go
	bin/bin

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/main.go