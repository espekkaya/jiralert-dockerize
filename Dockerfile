FROM golang:1.8

WORKDIR /go/src/app

RUN go get github.com/espekkaya/jiralert-dockerize/jiralert-dockerize/cmd/jiralert

COPY config/jiralert.yml /go/src/app/config/jiralert.yml
COPY config/jiralert.tmpl /go/src/app/config/jiralert.tmpl

ENTRYPOINT ["jiralert"]