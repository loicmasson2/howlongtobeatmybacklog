FROM golang:1.21.4

RUN go install github.com/cespare/reflex@latest

COPY reflex.conf /usr/local/etc/
COPY build.sh /usr/local/bin/

WORKDIR /app

VOLUME /go

CMD ["reflex", "-d", "none", "-c", "/usr/local/etc/reflex.conf"]

#WORKDIR /go/src/app
#
#COPY . .
#
#RUN go build -o main main.go
#
#CMD ["./main"]