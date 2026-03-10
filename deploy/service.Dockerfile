FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o person-enrichment-service ./cmd/service

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/person-enrichment-service /usr/local/bin/
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs/swagger /app/docs/swagger


ENTRYPOINT ["/usr/local/bin/person-enrichment-service"]