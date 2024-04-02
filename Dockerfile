# -- build stage --
FROM golang:1.21.5-alpine3.18 as build

WORKDIR /usr/src

RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./exporter ./


# -- runtime stage --
FROM scratch as runtime

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/src/exporter /exporter

ENTRYPOINT ["/exporter"]
