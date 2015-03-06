FROM golang:1.4.1-onbuild

export GOPATH=~/go
mkdir $GOPATH/{src,bin,pkg}