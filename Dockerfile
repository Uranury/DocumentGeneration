# ---- builder: build Go binary ----
FROM golang:1.23.1-alpine AS builder

# build deps
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main ./cmd/api/main.go

# ---- runtime: minimal Alpine ----
FROM alpine:latest

# install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates bash

WORKDIR /opt/app

# copy binary and templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
 && chown -R appuser:appgroup /opt/app

USER appuser

ENV PATH="/opt/app:${PATH}"
EXPOSE 8080

CMD ["./main"]
