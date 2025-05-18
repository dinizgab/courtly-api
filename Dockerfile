FROM golang:1.23.3-alpine

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/main ./cmd/main.go

CMD ["/app/bin/main"]
EXPOSE 8000
