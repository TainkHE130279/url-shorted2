# Dockerfile chính cho ứng dụng URL Shortener
FROM url-shortener:golang AS app

# Set working directory
WORKDIR /app

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY go.mod go.sum ./

# Build ứng dụng với optimizations
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o main ./cmd/main.go

# Tạo thư mục cho database
RUN mkdir -p /app/data

# Chuyển ownership cho appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run ứng dụng
CMD ["./main"]
