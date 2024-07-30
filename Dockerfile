# build stage
FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
# COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

# Ensure the start.sh script has execute permissions
RUN chmod +x start.sh wait-for.sh
RUN dos2unix start.sh wait-for.sh

EXPOSE 5001 50051
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]