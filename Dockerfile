# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /backendim-brain ./cmd/server

# Runtime stage
FROM alpine:3.19
WORKDIR /app

# Install core dependencies with verified versions
RUN apk add --no-cache \
  ca-certificates=20241121-r1 \
  curl=8.12.1-r0 \
  python3=3.11.11-r0 \
  py3-pip=23.3.1-r0 \
  git=2.43.6-r0 \
  bash=5.2.21-r0 \
  jq=1.7.1-r0

# Install security tools with verified versions
RUN apk add --no-cache --virtual .security-deps \
  openssl=3.1.4-r5 \
  libcrypto3=3.1.4-r5

# Install AWS CLI and kubectl
COPY scripts/install-awscli.sh scripts/install-kubectl.sh /tmp/
RUN /tmp/install-awscli.sh && \
  /tmp/install-kubectl.sh && \
  rm -f /tmp/install-*.sh && \
  rm -rf /var/cache/apk/*

# Application setup
COPY --from=builder /backendim-brain .
COPY scripts/ ./scripts/

# Security hardening
RUN find ./scripts/ -type f -name '*.sh' -exec chmod 0755 {} + && \
  adduser -D -u 1001 backenduser && \
  chown -R backenduser:backenduser /app && \
  chmod 0700 /root/.aws /root/.kube

ENV KUBECONFIG=/root/.kube/config \
  AWS_EC2_METADATA_DISABLED=true \
  PATH="/app/scripts:${PATH}" \
  GIT_SSL_NO_VERIFY="false"

USER backenduser

HEALTHCHECK --interval=30s --timeout=3s CMD scripts/healthcheck.sh
ENTRYPOINT ["/app/scripts/kube-init.sh"]
CMD ["./backendim-brain"]
