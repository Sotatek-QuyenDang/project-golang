FROM golang:1.24.2
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]

