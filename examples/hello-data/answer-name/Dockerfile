FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./main.go ./go.mod /app/
COPY ./vendor/ /app/vendor/

RUN go mod tidy
RUN go build -mod=vendor -o main main.go

FROM alpine:latest

COPY --from=builder /app/main /app/ 
WORKDIR /app
CMD ["./main"]