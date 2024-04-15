# coz we have 1.16 inside go.mod
FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# also add appropriate .dockerignore could be wise :)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o farawayweb .

# Possible to use scratch images, but skip it for now
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/farawayweb .
EXPOSE 8080
CMD ["./farawayweb"]
