FROM golang:1.24-alpine

WORKDIR /app/db-service

COPY db-service/go.mod db-service/go.sum ./

RUN go mod download

COPY ./db-service ./

RUN go build -o /app/bin/db-app ./cmd/app/main.go

CMD ["/app/bin/db-app"]

EXPOSE 44044
