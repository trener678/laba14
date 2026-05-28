FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/collector ./cmd/collector

FROM alpine:3.22
RUN adduser -D -H appuser
USER appuser
COPY --from=build /out/collector /collector
EXPOSE 8080
ENTRYPOINT ["/collector"]
