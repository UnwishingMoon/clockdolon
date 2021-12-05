#Build
FROM golang:1.17 AS build

WORKDIR /go-app

COPY * ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /app

# Deploy
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build /app /

USER nonroot:nonroot

ENTRYPOINT ["/app"]