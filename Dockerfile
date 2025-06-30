FROM alpine:3.20

# Install required runtime dependencies for CGO
RUN apk add --no-cache libc6-compat sqlite curl tar

WORKDIR /app

# Install lazysql
RUN curl -L https://github.com/jorgerojas26/lazysql/releases/download/v0.4.0/lazysql_Linux_x86_64.tar.gz \
  | tar -xz -C /usr/local/bin && \
  chmod +x /usr/local/bin/lazysql

# Copy binary and static assets from GitHub Actions build
COPY ./main /app/main
RUN mkdir /data


CMD ["/app/main", "serve"]
