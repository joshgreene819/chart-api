FROM golang:1.18.2-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./chart-api

EXPOSE 8080

# Want to move this to docker-compose.yaml eventually. More sensible since referencing the 'db'
# service name
ENV MONGO_URI="mongodb://db:27107" 

ENTRYPOINT [ "./chart-api" ]