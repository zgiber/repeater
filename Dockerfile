FROM golang:1.8-alpine

RUN mkdir -p /app
COPY . /app
WORKDIR /app
RUN go build
ENTRYPOINT /app/app
