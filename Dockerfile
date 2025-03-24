FROM ubuntu:22.04

# Install Go and necessary packages
RUN apt-get update && apt-get install -y \
    golang-go \
    software-properties-common && \
    add-apt-repository -y ppa:dqlite/dev && \
    apt-get update && \
    apt-get install -y dqlite-tools-v3 libdqlite-dev

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy application source
COPY . .

# Build the binary with the necessary tags
RUN go build -tags libsqlite3 -o dqlite-demo ./main.go
RUN mkdir -p /dqlite

# Expose the API port
EXPOSE 8001

ENTRYPOINT ["./dqlite-demo"]
CMD ["--api", "0.0.0.0:8001", "--dir", "/dqlite", "--db", "0.0.0.0:9001"]
