# ---- Build stage ----
FROM golang:1.21-alpine3.18 AS builder

RUN apk update && apk upgrade --available && sync

WORKDIR /app
COPY . .

# Make sure ./cmd/fsb/main.go exists
RUN CGO_ENABLED=0 go build -o /app/fsb -ldflags="-w -s" ./cmd/fsb

# ---- Final stage ----
FROM scratch

COPY --from=builder /app/fsb /app/fsb

EXPOSE ${PORT}

ENTRYPOINT ["/app/fsb", "run"]
