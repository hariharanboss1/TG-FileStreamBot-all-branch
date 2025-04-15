FROM golang:1.21-alpine3.18 as builder

# Install build dependencies
RUN apk update && apk upgrade --available && \
    apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /app/fsb -ldflags="-w -s" ./cmd/fsb

FROM scratch
COPY --from=builder /app/fsb /app/fsb
EXPOSE ${PORT}
ENTRYPOINT ["/app/fsb", "run"]
