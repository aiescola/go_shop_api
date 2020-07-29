FROM golang:1.14.6-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY server/go.mod .
COPY server/go.sum .

RUN go mod download 

COPY server/ .

RUN go build

FROM alpine:latest

ENV BBDD_URI=$BBDD_URI
ENV BBDD_NAME=$BBDD_NAME
ENV SESSION_KEY=$SESSION_KEY
ENV PORT=$PORT

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app .

CMD ["./go_shop_api"]



