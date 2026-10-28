ARG NODE_VERSION=22
ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.23

FROM node:${NODE_VERSION}-alpine AS web
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY web/ ./
RUN npm run build

FROM golang:${GO_VERSION}-alpine AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY migrations/ ./migrations/
RUN CGO_ENABLED=0 go build -trimpath -o /out/chronicle ./cmd/chronicle

FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 10001 chronicle \
    && adduser -D -u 10001 -G chronicle chronicle \
    && mkdir -p /data /etc/chronicle \
    && chown chronicle:chronicle /data
WORKDIR /app
COPY --from=backend /out/chronicle /usr/local/bin/chronicle
COPY --from=web /src/web/build ./web/build
COPY config/docker.json /etc/chronicle/config.json
ENV CHRONICLE_CONFIG=/etc/chronicle/config.json
USER chronicle
EXPOSE 8080
VOLUME /data
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["chronicle", "healthcheck"]
ENTRYPOINT ["chronicle"]
