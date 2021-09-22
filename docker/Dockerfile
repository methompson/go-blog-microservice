FROM alpine:3.14
# FROM golang:1.16

RUN mkdir -p /srv/blog
WORKDIR /srv/blog

COPY bin/blog-microservice ./blog-microservice
COPY firebase.json ./firebase.json

CMD ["./blog-microservice"]