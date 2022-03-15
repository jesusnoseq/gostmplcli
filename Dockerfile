# syntax=docker/dockerfile:1

FROM golang:1.17 AS builder
WORKDIR /github.com/jesusnoseq/gostmplcli

COPY go.mod *.go ./

RUN CGO_ENABLED=0 go build -o gostmplcli


FROM gcr.io/distroless/static
WORKDIR /app

COPY --from=builder /github.com/jesusnoseq/gostmplcli/gostmplcli /

ENTRYPOINT ["/gostmplcli"]