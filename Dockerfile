# Base image
FROM golang:1.24 AS builder

WORKDIR /app

ADD . /app/

RUN go build -o /app/leanbot .

# Trim image
FROM golang:1.24 AS runner

WORKDIR /app

COPY --from=builder /app/leanbot /app/leanbot

ENTRYPOINT ["/app/leanbot"]
