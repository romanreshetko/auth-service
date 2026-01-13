FROM golang:1.24-alpine

WORKDIR /app
RUN mkdir -p /app/keys
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service

EXPOSE 8080
CMD ["./auth-service"]