FROM golang:1.4

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]

COPY . /go/src/app
RUN rm tutum_decoder.go

RUN go-wrapper download
RUN go test
