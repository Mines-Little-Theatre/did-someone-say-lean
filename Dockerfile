# Base image
FROM golang:1.24-alpine AS builder

RUN apk add git

WORKDIR /go/src/app

ADD . .

RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

# Trim image
# FROM golang:1.24 AS runner
FROM scratch AS runner

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/did-someone-say-lean /leanbot

ENTRYPOINT ["/leanbot"]
