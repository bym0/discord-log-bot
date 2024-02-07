# Build Stage
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

# App Stage
FROM alpine:latest  
WORKDIR /root/
COPY --from=build /app/app .

CMD ["./app"]
