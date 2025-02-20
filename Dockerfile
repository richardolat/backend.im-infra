# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /backendim-brain ./cmd/server

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache \
    ca-certificates \
    aws-cli \
    curl \
    python3 \
    py3-pip \
    git && \
    pip3 install --upgrade awscli && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    mkdir -p /root/.aws /root/.kube

WORKDIR /app
COPY --from=builder /backendim-brain .
COPY scripts/ ./scripts/
RUN chmod +x ./scripts/*.sh && \
    adduser -D -u 1001 backenduser && \
    chown -R backenduser:backenduser /app /root/.aws /root/.kube

ENV KUBECONFIG=/root/.kube/config \
    AWS_EC2_METADATA_DISABLED=true \
    PATH="/app/scripts:${PATH}"

USER backenduser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost:8080/ws || exit 1
CMD ["./backendim-brain"]
