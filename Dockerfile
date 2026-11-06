FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP
RUN go build -o /app/service ./cmd/${APP}

FROM alpine:3.22

COPY --from=builder /app/service /app/service

ENTRYPOINT ["/app/service"]