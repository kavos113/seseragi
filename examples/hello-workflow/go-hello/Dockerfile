FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./main.go /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:latest

COPY --from=builder /app/main /app/ 
WORKDIR /app
CMD ["./main"]