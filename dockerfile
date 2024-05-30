FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build main.go

EXPOSE 1323

ENTRYPOINT /app/main