FROM golang:1.25-trixie AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /exporter

FROM gcr.io/distroless/static-debian12:latest AS run

COPY --from=build /exporter /exporter

ENTRYPOINT ["/exporter", "start"]