FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8000


CMD ["/docker-gs-ping"]