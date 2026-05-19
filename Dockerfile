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
      /app/public/uploads/payments \
      /app/public/uploads/payments_pelunasan \
      /app/public/uploads/teams/id_cards \
      /app/public/uploads/teams/logos \
      /app/public/uploads/teams/photos \
      /app/public/uploads/teams/rekomendasi \
      /app/public/uploads/wallets && \
    chown -R app:app /app

COPY --from=builder /out/paskihub-be /app/paskihub-be

USER app

EXPOSE 3010

CMD ["/app/paskihub-be"]
