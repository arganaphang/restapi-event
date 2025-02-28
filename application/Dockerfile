FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/application ./main.go

FROM alpine
WORKDIR /usr/bin
ENV DATABASE_HOST "${DATABASE_HOST}"
ENV DATABASE_NAME "${DATABASE_NAME}"
ENV DATABASE_USER "${DATABASE_USER}"
ENV DATABASE_PASSWORD "${DATABASE_PASSWORD}"

EXPOSE ${DATABASE_HOST}
EXPOSE ${DATABASE_NAME}
EXPOSE ${DATABASE_USER}
EXPOSE ${DATABASE_PASSWORD}
EXPOSE 8000

COPY --from=builder /app/application /usr/bin/application
CMD ["application"]