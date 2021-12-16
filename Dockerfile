#Build
FROM --platform=$BUILDPLATFORM golang:1.17 AS build

WORKDIR /app

COPY . .
RUN go mod download && \
    CGO_ENABLED=0 go build -o /app/app -ldflags "-s -w" cmd/main.go

# Deploy
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build /app/app /app/config.toml /

USER nonroot:nonroot

ENTRYPOINT ["/app"]