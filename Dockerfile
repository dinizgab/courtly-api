FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 \
    GOFLAGS="-trimpath -mod=readonly"

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /bin/app ./cmd/main.go

FROM alpine:3.20  AS runtime

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app
COPY --from=builder /bin/app /app/app

ENV PORT=8000
EXPOSE 8000

RUN adduser -D -u 10001 appuser
USER appuser

ENTRYPOINT ["/app/app"]

