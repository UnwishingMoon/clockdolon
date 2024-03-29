#Build
FROM --platform=$BUILDPLATFORM golang:latest AS build

WORKDIR /app
COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/app -ldflags "-s -w" cmd/main.go

# Deploy
FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian11

WORKDIR /
COPY --from=build /app/app /

USER nonroot:nonroot

ENTRYPOINT ["/app"]