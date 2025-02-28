set dotenv-load

setup:
	go install github.com/bitfield/gotestdox/cmd/gotestdox@latest

format:
	gofumpt -l -w .
	goimports-reviser -rm-unused -set-alias ./...
	golines -w -m 120 .

# health -> Hit Health Check Endpoint
health:
	curl -s http://localhost:8000/healthz | jq

# migrate-create -> create migration
migrate-create NAME:
	migrate create -ext sql -dir ./migrations -seq {{NAME}}

# migrate-up -> up migration
migrate-up:
	migrate -path ./migrations -database "$DATABASE_URL" up

# migrate-down -> down migration
migrate-down:
	migrate -path ./migrations -database "$DATABASE_URL" down


# build -> build application
build APP:
	go build -o ./dist/{{APP}} ./{{APP}}/cmd/main.go

# run -> application
run APP:
	./dist/{{APP}}

# dev -> run build then run it
dev APP: 
	watchexec -r -c -e go -- just build run {{APP}}

# test -> testing
test APP:
  gotestdox -v ./{{APP}}/...

test-load:
  hey -n 100000 -c 50 -m POST -H "Content-Type: application/json" -D ./request_body.txt http://localhost:8000/save