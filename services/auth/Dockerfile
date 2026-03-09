FROM golang:alpine3.23 AS builder

RUN apk add --no-cache git

WORKDIR /workspace

COPY . .

ENV GOWORK=/workspace/go.work

WORKDIR /workspace/services/auth

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o app ./cmd/

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /workspace/services/auth/app /app

USER nonroot:nonroot

EXPOSE 50052

ENTRYPOINT ["/app"]