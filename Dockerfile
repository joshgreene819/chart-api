FROM golang:1.18.2-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./chart-api

EXPOSE 8080

ENTRYPOINT [ "./chart-api" ]