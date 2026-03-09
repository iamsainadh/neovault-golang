FROM golang:alpine3.23 AS builder

RUN apk add --no-cache git

WORKDIR /workspace

COPY . .

ENV GOWORK=/workspace/go.work

RUN go work sync

WORKDIR /workspace/services/api-gateway

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o app ./cmd/

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /workspace/services/api-gateway/app /app

USER nonroot:nonroot

EXPOSE 8081

ENTRYPOINT ["/app"]