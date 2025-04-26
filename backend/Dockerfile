FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM alpine:latest
COPY --from=builder /api /api
CMD ["/api"]