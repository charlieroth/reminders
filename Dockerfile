# Build the Go Binary.
FROM golang:1.23 AS build_reminders
WORKDIR /app
COPY . .
COPY ./vendor ./vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Run the Go Binary in Alpine.
FROM alpine:3.20
COPY --from=build_reminders /app/server /app
CMD ["./server"]