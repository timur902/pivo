FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/order-service ./order-service/cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/order-service /order-service
EXPOSE 50051
USER nonroot:nonroot
ENTRYPOINT ["/order-service"]
