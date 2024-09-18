FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN go build -o /payment-worker .

FROM alpine:3.20
COPY --from=build /payment-worker /payment-worker
ENTRYPOINT ["/payment-worker"]
