FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
COPY go.mod go.sum ./
RUN go mod tidy
RUN CGO_ENABLE=0 go build -o batcher_http ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/batcher_http .
COPY --from=builder /app/Configs/config.json ./Configs/config.json
EXPOSE 4000
CMD ["./batcher_http"]