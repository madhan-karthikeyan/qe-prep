FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build ./...

FROM golang:1.23-alpine
WORKDIR /app
COPY --from=builder /app .
CMD ["go", "test", "-race", "-count=1", "./..."]
