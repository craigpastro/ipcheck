FROM golang:1.16.5-buster AS builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY ./ .

RUN go mod download

RUN go build -o main .

FROM golang:1.16.5-buster

EXPOSE 8080

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]
