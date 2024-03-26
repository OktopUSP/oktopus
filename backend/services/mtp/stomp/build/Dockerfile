FROM golang:1.22@sha256:82e07063a1ac3ee59e6f38b1222e32ce88469e4431ff6496cc40fb9a0fc18229 as builder
WORKDIR /app
COPY ../ .
RUN CGO_ENABLED=0 GOOS=linux go build -o stomp cmd/stomp/main.go

FROM alpine:3.14@sha256:0f2d5c38dd7a4f4f733e688e3a6733cb5ab1ac6e3cb4603a5dd564e5bfb80eed
COPY --from=builder /app/stomp /
ENTRYPOINT ["/stomp"]