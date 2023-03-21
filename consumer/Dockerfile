FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /usr/goapplication ./src/main.go

FROM alpine
WORKDIR /usr/bin
COPY --from=builder /usr/goapplication /usr/bin/goapplication
CMD ["goapplication"]