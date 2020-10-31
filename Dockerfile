FROM golang:latest

COPY . /var/www

WORKDIR /var/www

RUN go get
    
RUN go build

CMD ["./golang-api"]