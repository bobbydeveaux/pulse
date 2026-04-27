# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:latest AS builder

RUN apk add --no-cache git

WORKDIR /src
COPY app/ .
RUN go build -o /pulse ./cmd/pulse

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache git bash

COPY --from=builder /pulse /usr/local/bin/pulse
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
