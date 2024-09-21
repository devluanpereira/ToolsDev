# Dockerfile

FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/public ./public
COPY --from=builder /app/static ./static
COPY --from=builder /app/web ./web

EXPOSE 8000

CMD ["./main"]