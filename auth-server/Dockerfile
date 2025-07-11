FROM golang:1.24.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go mod download

RUN ls -l /app/cmd/config.yaml

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]
