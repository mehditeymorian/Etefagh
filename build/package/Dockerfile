FROM golang:alpine AS builder

WORKDIR /app


COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /cmd

FROM alpine:latest

LABEL maintainer="Mehdi Teymorian <mehditeymorian322@gmail.com>"

WORKDIR /app/

COPY --from=builder /cmd .

EXPOSE 3000

ENTRYPOINT ["./cmd"]

CMD ["server"]