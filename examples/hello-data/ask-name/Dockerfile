FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./main.go ./go.mod /app/

RUN go mod tidy
RUN go build -o main main.go

FROM alpine:latest

COPY --from=builder /app/main /app/ 
WORKDIR /app
CMD ["./main"]