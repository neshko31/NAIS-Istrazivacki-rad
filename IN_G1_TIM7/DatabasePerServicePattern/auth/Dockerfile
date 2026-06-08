FROM golang:1.25-alpine AS builder

WORKDIR /auth_service

COPY auth_service/go.mod auth_service/go.sum ./
RUN go mod download

COPY auth_service/ .
RUN go build -o main .

FROM alpine:latest

WORKDIR /auth_service
COPY --from=builder /auth_service/main .

EXPOSE 8082

CMD ["./main"]

