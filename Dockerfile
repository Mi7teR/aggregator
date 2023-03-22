FROM golang:1.20.1 as builder

WORKDIR /go/src/app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go mod download

COPY . /go/src/app
RUN go build -o server cmd/main.go


FROM alpine:latest
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/server server
USER 2000
ENV TIMEOUT=30s
ENV HTTP_PORT=8080
ENTRYPOINT ./server