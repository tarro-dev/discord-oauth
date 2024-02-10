FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/web ./cmd/web

FROM gcr.io/distroless/base

COPY --from=builder /go/bin/web /
COPY static /static

CMD ["/web"]