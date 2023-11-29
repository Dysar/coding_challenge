# Coding challenge

I chose RESTful API. I added a UI as well. 
And you can configure the pack sizes without having to 
change the code as they are located in the `conf.json` configuration file.

## Build & Run

Running with Docker
```bash
docker build -t challenge .
docker run -it -p 8080:8080 --rm challenge
```

Running with GNU Make
```bash
make dev
```

Running with Go:
```Bash
go build -o bin/cmd cmd/main.go
bin/cmd
```

## How to use

* UI: https://stark-river-82961-278b69188afd.herokuapp.com

![UI](UI2.png)
* In case you would like to use the API, sample curls: 
```bash 
curl --location --request POST 'https://stark-river-82961-278b69188afd.herokuapp.com/api/v1/calculate_packs' \
--header 'Content-Type: application/json' \
--data-raw '{
"order_quantity":10
}'
```

Change the packs
```bash 
curl --location --request PUT 'https://stark-river-82961-278b69188afd.herokuapp.com/api/v1/pack_sizes' \
--header 'Content-Type: application/json' \
--data '{
    "pack_sizes":[23,31,53]
}'
```

Get the packs
```bash 
curl --location --request GET 'https://stark-river-82961-278b69188afd.herokuapp.com/api/v1/pack_sizes' 
```

### Dependency Versions
* GNU Make 3.81 (optional)
* Go 1.21