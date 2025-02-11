# Build stage
FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.20 
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "serve" ]
ENTRYPOINT [ "/app/main" ]
