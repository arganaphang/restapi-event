set dotenv-load

setup:
	go install github.com/bitfield/gotestdox/cmd/gotestdox@latest

# format -> format code
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
	echo "${DATABASE_PASSWORD}"
	migrate -path ./migrations -database "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@localhost:5432/$DATABASE_NAME?sslmode=disable" up

# migrate-down -> down migration
migrate-down:
	migrate -path ./migrations -database "$DATABASE_URL" down


# build -> build application
build APP:
	go build -o ./dist/{{APP}} ./{{APP}}/main.go
# run -> application
run APP:
	./dist/{{APP}}

# dev -> run build then run it
dev APP: 
	watchexec -r -c -e go -- just build {{APP}} run {{APP}}

# test -> testing
test APP:
  gotestdox -v ./{{APP}}/...

test-load N="100000" C="50":
  hey -n {{N}} -c {{C}} -m POST -H "Content-Type: application/json" -D ./request_body.txt http://localhost:8000/save