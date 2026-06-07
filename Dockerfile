FROM golang:1.24.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/paskihub-be ./cmd/app/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && \
    adduser -S app -G app && \
    mkdir -p \
      /app/data/logs \
      /app/public/uploads/events \
      /app/public/uploads/teams/logos \
      /app/storage/private/payments \
      /app/storage/private/payments_pelunasan \
      /app/storage/private/teams/id_cards \
      /app/storage/private/teams/photos \
      /app/storage/private/teams/rekomendasi \
      /app/storage/private/wallets && \
    chown -R app:app /app

COPY --from=builder /out/paskihub-be /app/paskihub-be

USER app

EXPOSE 3010

CMD ["/app/paskihub-be"]
