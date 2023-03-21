FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /usr/goapplication ./src/main.go

FROM alpine
WORKDIR /usr/bin
COPY --from=builder /usr/goapplication /usr/bin/goapplication
EXPOSE 8000
CMD ["goapplication"]