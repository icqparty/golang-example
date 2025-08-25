FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN apk add --no-cache bash
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/main ./main
RUN chmod +x /app/main
EXPOSE 80 443
CMD ["/app/main","-f","./config/config.yaml"]