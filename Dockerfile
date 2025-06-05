FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o dns-server main.go
EXPOSE 53/udp
CMD ["./dns-server"]
