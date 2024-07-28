# build stage
FROM golang:1.20-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
# COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 5000
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]