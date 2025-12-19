FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api-gateway ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/api-gateway /api-gateway
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api-gateway"]
