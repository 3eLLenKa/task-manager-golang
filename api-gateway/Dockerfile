FROM golang:1.24-alpine

WORKDIR /app/api-gateway

COPY api-gateway/go.mod api-gateway/go.sum ./

RUN go mod download

COPY api-gateway ./

RUN go build -o /app/bin/api-app ./cmd/app/main.go

CMD ["/app/bin/api-app"]

EXPOSE 8080
