# Coding challenge

I chose RESTful API. I added a UI as well. 
And you can configure the pack sizes without having to 
change the code as they are located in the `conf.json` configuration file.

## Build & Run
Running with GNU Make
```bash
make dev
```

Running with Go:
```Bash
go build -o bin/bin cmd/main.go
bin/bin
```

## How to use

* UI: https://stark-river-82961-278b69188afd.herokuapp.com

![UI](UI2.png)
* In case you would like to use the API, sample curl: 
```bash 
curl --location --request POST 'https://stark-river-82961-278b69188afd.herokuapp.com/api/v1/calculate_packs' \
--header 'Content-Type: application/json' \
--data-raw '{
"order_quantity":10
}'
```

### Dependency Versions
* GNU Make 3.81 (optional)
* Go 1.21