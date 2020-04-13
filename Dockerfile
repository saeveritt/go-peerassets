FROM golang:1.14


ENV GOOS=linux  \
    GOARCH=amd64 \
    APP_ENV=docker

COPY ./docker-entrypoint.sh  /usr/local/bin/docker-entrypoint.sh
COPY ./app/config/walletnotify.sh /usr/local/bin/walletnotify.sh

RUN chmod +x  /usr/local/bin/docker-entrypoint.sh
RUN chmod +x  /usr/local/bin/walletnotify.sh

ENTRYPOINT ["docker-entrypoint.sh"]

COPY ./app /app

WORKDIR /app

RUN go mod download

RUN go build -o main .

CMD ["/app/main"]
ENV LISTEN_PORT 8089
ENV LISTEN_PORT 9904
EXPOSE 8089