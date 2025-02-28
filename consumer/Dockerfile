FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /usr/goapplication ./src/main.go

FROM alpine
WORKDIR /usr/bin
ENV BROKER_URLS "${BROKER_URLS}"
ENV DATABASE_HOST "${DATABASE_HOST}"
ENV DATABASE_NAME "${DATABASE_NAME}"
ENV DATABASE_USER "${DATABASE_USER}"
ENV DATABASE_PASSWORD "${DATABASE_PASSWORD}"

EXPOSE ${BROKER_URLS}
EXPOSE ${DATABASE_HOST}
EXPOSE ${DATABASE_NAME}
EXPOSE ${DATABASE_USER}
EXPOSE ${DATABASE_PASSWORD}

COPY --from=builder /usr/goapplication /usr/bin/goapplication
CMD ["goapplication"]