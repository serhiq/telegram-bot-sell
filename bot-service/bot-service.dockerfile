# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN #CGO_ENABLED=0 go build -o bot ./
RUN CGO_ENABLED=0 go build -o bot ./cmd/sbot/


RUN chmod +x /app/bot

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/bot /app

#CMD [ "/app/bot" ]
CMD [ "/app/bot" ]
