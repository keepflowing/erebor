FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o backend


# final stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["./backend"]
