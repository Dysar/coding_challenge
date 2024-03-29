# Coding challenge

I chose RESTful API. I added a UI as well. The user is able to get an overview the pack sizes,
modify them and calculate the required amount of packs


## The algorithm 

Requirement: 
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out no more items than necessary to fulfil the order.
3. Within the constraints above, send out Rules 1 & 2 send out as few packs as possible to fulfil each order.

Solution:
1. The pack sizes are sorted in descending order
2. We iterate over the pack sizes and divide the ordered quantity by the pack size to find out the quotient and the remainder 
3. Next, using the quotient we add a number of packs to the result and using the remainder decide whether we can stop iterating immediately (remainder is less than the small pack size), or we need to iterate further
4. After the initial iteration is done, we check if there are any other combinations that would suffice the initial equation x1n + x2m ... x3p = total sum; where n,m,p are the pack sizes, x1, x2, x3 are the quantities of the packs. So that x1 + x2 ... xn = Minimize Z 
5. After that, we try to replace packs of the same size with bigger-sized packs to optimize the pack count

## Logging

For this test assignment, I decided to add excessive logging to the algorithm to showcase its work principles.


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

* UI: https://coding-challenge-app-2f0da67d681f.herokuapp.com

![UI](ui.png)
* In case you would like to use the API, sample curls: 
```bash 
curl --location --request POST 'https://coding-challenge-app-2f0da67d681f.herokuapp.com/api/v1/calculate_packs' \
--header 'Content-Type: application/json' \
--data-raw '{
"order_quantity":10
}'
```

Change the packs
```bash 
curl --location --request PUT 'https://coding-challenge-app-2f0da67d681f.herokuapp.com/api/v1/pack_sizes' \
--header 'Content-Type: application/json' \
--data '{
    "pack_sizes":[23,31,53]
}'
```

Get the packs
```bash 
curl --location --request GET 'https://coding-challenge-app-2f0da67d681f.herokuapp.com/api/v1/pack_sizes' 
```

### Dependency Versions
* GNU Make 3.81 (optional)
* Go 1.21