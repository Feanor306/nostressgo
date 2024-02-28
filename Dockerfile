FROM golang:1.21.4-alpine3.18
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o ./bin/main ./cmd/server/main.go

EXPOSE 3000

CMD ["/app/bin/main"]