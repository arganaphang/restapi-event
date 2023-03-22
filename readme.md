# Simple RESTfull API

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
- [Tasks](https://taskfile.dev/) Runner (Makefile alternative)
- [Hey](https://github.com/rakyll/hey) Http load Testing

## How to run?

### Using Task runner

1. to start the postgres, redpanda, consumer, and application
```sh
task start
```
you can access [Redpanda Console](http://localhost:8080/overview) to see redpanda messages
2. to run load testing
```sh
task test:load
```
this command will seed the traffic based on `n` request at `Taskile.yaml`, i set it 100000 requests
the request body you can see it at `request_body.txt` for data i give it 2 transaction per request, so it will insert n * 2 into postgres

### Manual

All of the commands that i using already documented at `Taskfile.yaml`
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
> We can see the average for each request is <ins>*0.0077 secs*</ins>, fastest is <ins>*0.0004 secs*</ins>, slowest is <ins>*0.1533 secs*</ins>

![Redpanda Console](/assets/redpanda_console.png)
> <ins>200,000 Messages</ins>

![Database Count](/assets/database_1.png)
> select count(id) from transactions;

![Database Count](/assets/database_2.png)
> select * from transactions limit 300; -- of course we need to limit this query to prevent database client crash
> field insert wrong field for `price` and already fixed in [This Commit](/commit/a416fb36d444143e8771c4164ee44ccae7921498#diff-c148edd2cbfefa32169ce5dbef35c4e81413f7212a8140d69f4242689baa2321R154)

## Architecture

![Architecture Diagram](/assets/architecture.png)

## Todos

- [ ] Optimize consumer using worker pool