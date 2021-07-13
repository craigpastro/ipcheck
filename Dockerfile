FROM golang:1.16-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV GOOS=linux

COPY ./ .

RUN go mod download

RUN go build -o main .

FROM golang:1.16-alpine

EXPOSE 50051

WORKDIR /app

COPY ./.env .
COPY --from=builder /app/main .

CMD ["/app/main"]
