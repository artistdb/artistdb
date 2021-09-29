FROM golang:1.17.1-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY *.go ./

RUN go build -o artistdb

FROM alpine

COPY --from=builder /build/artistdb /api/artistdb

ENTRYPOINT ["/api/artistdb"]