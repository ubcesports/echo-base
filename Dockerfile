FROM golang:1.25 AS builder
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /go/src/app/cmd/server
RUN CGO_ENABLED=0 go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian12 AS runtime
COPY --from=builder /go/bin/app /usr/local/bin/echobase
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/echobase"]
