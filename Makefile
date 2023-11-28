dev:
	go mod tidy
	go build -o app cmd/main.go
	app