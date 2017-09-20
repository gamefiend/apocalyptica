FROM golang:1.8
STOPSIGNAL SIGTERM
LABEL maintainer="Quinn Murphy"

WORKDIR /go/src/github.com/gamefiend/apocalyptica
  
COPY . .

RUN curl https://glide.sh/get | sh \
    && glide install \
    && glide up \
    && CGO_ENABLED=0 go-wrapper install
EXPOSE 8080
ENTRYPOINT ["go-wrapper","run"]
