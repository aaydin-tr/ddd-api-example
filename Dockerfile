FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install go.uber.org/mock/mockgen@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go generate ./...
RUN go test ./... 
RUN swag init -g ./controller/ticket/ticket.go --parseDependency true

RUN go build -o main ./cmd/main.go


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env .env

EXPOSE ${PORT}

CMD ["./main"]