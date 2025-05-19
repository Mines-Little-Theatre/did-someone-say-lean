# Base image
FROM golang:1.24 AS builder

WORKDIR /app

ADD . /app/

RUN go build -o /leanbot .

# Trim image
FROM golang:1.24 AS runner

WORKDIR /app

COPY --from=builder /leanbot /app/leanbot

CMD ["/app/leanbot"]
