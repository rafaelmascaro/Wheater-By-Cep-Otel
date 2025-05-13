FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/microservice

FROM alpine:latest
COPY --from=builder /app/cmd/microservice/.env .
COPY --from=builder /app/server .
CMD ["./server"]