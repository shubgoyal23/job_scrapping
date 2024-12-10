FROM golang:1.23.3-alpine

WORKDIR /app

WORKDIR /app

# Copy the application files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

CMD ["./app"]