FROM golang:1.11

ENV APP_DIR $GOPATH/src/github.com/aereal/poc-nginx-healthcheck
RUN mkdir -p ${APP_DIR}
ADD . ${APP_DIR}
WORKDIR ${APP_DIR}
RUN go build -o ./web.app ./app
ENTRYPOINT [ "./web.app" ]
