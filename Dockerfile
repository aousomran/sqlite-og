FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o sqliteogd -a -ldflags '-w -extldflags "-static"' ./cmd/sqliteogd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

RUN apk --no-cache add sqlite

CMD ["/app/sqliteogd"]
