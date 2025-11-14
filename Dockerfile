FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
# RUN go mod download
COPY . .

RUN go build -mod=vendor -o main ./cmd/api-gateway # -o main = output compiled file as main

CMD ["./main"]

# CMD ["sh", "-c", "go run ./cmd/api-gateway"]