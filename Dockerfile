FROM golang:1.14


ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux  \
    GOARCH=amd64 \

WORKDIR /build

COPY app/go.mod .
COPY app/go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main

EXPOSE 8089

CMD["/dist/main"]