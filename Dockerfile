FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
