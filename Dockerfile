FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN templ generate
RUN npx tailwindcss -i ./input.css -o ./static/output.css --minify

RUN go build -o app .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./app"]