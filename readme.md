# Simple RESTfull API (Microservice)

Simple Service using RESTFullAPI and Event Driven

## Tools
> for the service already using container, so you just need to pull the docker image

1. Golang v1.20.1
  - Gin Gonic
  - SQLX + pq (postgres sql lib for golang)
  - Sarama (kafka library) *yes redpanda can be use using kafka client
2. PostgreSQL
3. Redpanda
4. Docker

## Prerequirement

- Docker + Docker Compose
- [Justfile](https://github.com/casey/just) Runner (Makefile alternative)
- [Watchexec](https://github.com/watchexec/watchexec) Watcher
- [Hey](https://github.com/rakyll/hey) Http load Testing

## How to run?

### Using Just runner

1. to start the postgres, redpanda, consumer, and application
```sh
docker-compose up -d --build
```
you can access [Redpanda Console](http://localhost:8080/overview) to see redpanda messages
2. to run load testing
```sh
just test-load
```
the request body you can see it at `request_body.txt` for data i give it 2 transaction per request, so it will insert n * 2 into postgres

### Manual

All of the commands that i using already documented at `Justfile`
See detail about this command above

1. Run Services
```sh
docker-compose up -d --build
```

2. Run Load Testing
```sh
hey -n 100000 -c 50 -m POST -H "Content-Type: application/json" -D ./request_body.txt http://localhost:8000/save
```

## Result

I've tested the server with 100000 Request, that's mean 2 * 100000 = 200000 data will be inserted into database

![Load Testing](/assets/load_test.png)
> We can see the average time to make all request is <ins>*5.7172 secs*</ins>, for each request is <ins>*0.0026 secs*</ins>, fastest is <ins>*0.0001 secs*</ins>, slowest is <ins>*0.6050 secs*</ins>, and Request/Sec is <ins>17491.2106</ins>

![Redpanda Console](/assets/redpanda_console.png)
> <ins>200,000 Messages</ins>

![Database Count](/assets/database_1.png)
> select count(id) from transactions;

![Database Count](/assets/database_2.png)
> select * from transactions limit 100; -- of course we need to limit this query to prevent database client crash

## Architecture

![Architecture Diagram](/assets/architecture.png)

## Todos

- [ ] Optimize consumer using worker pool