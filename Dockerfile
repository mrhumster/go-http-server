FROM golang:1.25-alpine AS builder
ARG VERSION=1.0.1 
ARG BUILD_DATE=17.09.2025 

WORKDIR /app
COPY go.mod ./ 
# Copy go.sum only if it exists
RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -X main.version=$VERSION -X main.buildDate=$BUILD_DATE" \
  -o server .

FROM alpine:3.18
ARG VERSION=unknown
ARG BUILD_DATE=unknown
LABEL version=$VERSION \
  build-date=$BUILD_DATE \
  maintainer="me@xomrkob.ru"
# for HTTPS
# RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 appgroup && \
  adduser -D -u 1000 -G appgroup appuser
WORKDIR /app 
COPY --from=builder --chown=appuser:appgroup /app/server .
USER appuser
EXPOSE 3333
EXPOSE 4444
CMD ["./server"]

