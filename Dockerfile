FROM golang:1.20-alpine3.17 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=arm64 go build -v -o tg-home-bot

FROM alpine:3.17

WORKDIR /usr/local/bin/

COPY --from=builder /usr/src/app/tg-home-bot /usr/src/app/.env ./

ENTRYPOINT ["/usr/local/bin/tg-home-bot"]