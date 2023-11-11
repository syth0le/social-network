FROM golang:1.21.4-alpine as builder

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /usr/local/bin/app cmd/social-network/main.go --config=local_config.yaml

FROM alpine:3.18

COPY --from=builder /usr/local/bin/app /app

CMD ["/app"]