FROM golang:bullseye as builder

RUN mkdir -p /app

COPY ./ /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontEnd ./cmd/web

RUN chmod +x /app/frontEnd

CMD ["/app/frontEnd"]