FROM golang:1.17-alpine

WORKDIR /app

COPY . /app

RUN go build -o main .

# Expose port 8080
EXPOSE 8080

CMD ["./main"]