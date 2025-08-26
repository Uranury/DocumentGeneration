# ---- builder: build Go binary on Debian-based image ----
FROM golang:1.23.1-bullseye AS builder

# build deps
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    git \
    tzdata \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# static-ish build; adjust if you rely on cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/api/main.go


# ---- runtime: Debian Bookworm slim with full LibreOffice + JRE ----
FROM debian:bookworm-slim

# install reliable libreoffice and OpenJDK (headless), plus fonts
RUN apt-get update && apt-get install -y \
        libreoffice \
        libreoffice-writer \
        libreoffice-java-common \
        openjdk-17-jre-headless \
        fonts-dejavu-core \
        fontconfig \
        ca-certificates \
        bash \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /opt/app

# copy binary + templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# create non-root user
RUN groupadd -g 1001 appgroup \
 && useradd -u 1001 -g appgroup -m appuser \
 && chown -R appuser:appgroup /opt/app

USER appuser

ENV PATH="/opt/app:${PATH}"
EXPOSE 8080

CMD ["./main"]
