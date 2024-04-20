FROM golang:1.22.2-alpine3.19

WORKDIR /BoomberLoad/cmd


COPY . /BoomberLoad

RUN go build -o main .

CMD ["./main"]