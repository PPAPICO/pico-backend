FROM golang:1.26-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main cmd/main.go

CMD ["/app/main"]