FROM golang:1.24.5-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o session-auth ./cmd/main.go

CMD ["/session-auth"]
