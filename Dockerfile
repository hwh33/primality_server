FROM golang:1.4.1-onbuild

RUN go get github.com/hwh33/primality_server

EXPOSE 8080
