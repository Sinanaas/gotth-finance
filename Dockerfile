FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.2.793
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate

RUN go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY app.env .

EXPOSE 8081
CMD ["./app"]