FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/notification-service ./notification-service/cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/notification-service /notification-service
USER nonroot:nonroot
ENTRYPOINT ["/notification-service"]
