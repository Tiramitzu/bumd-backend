FROM golang:1.20.7-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary -tags musl

ENTRYPOINT ["/app/binary"]