FROM golang:1.24.2
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .
EXPOSE 3000
CMD ["./main"]

