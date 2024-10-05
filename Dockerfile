FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y netcat-traditional

COPY . /app
RUN chmod +x /app/wait-for-postgres.sh

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/app/wait-for-postgres.sh", "/docker-gs-ping"]
