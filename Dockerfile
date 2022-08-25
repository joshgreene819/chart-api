FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./chart-api

EXPOSE 8080

CMD [ "./chart-api" ]