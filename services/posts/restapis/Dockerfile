FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /posts-api \
    ./cmd/posts-api

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /posts-api /posts-api

EXPOSE 8200

ENTRYPOINT ["/posts-api"]
CMD ["start"]