FROM golang:1.22.3-alpine3.19 as builder
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
WORKDIR /BoomberLoad/cmd
COPY . /BoomberLoad
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o boomberload


FROM scratch
WORKDIR /BoomberLoad/cmd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /BoomberLoad/cmd/boomberload .
#COPY --from=builder /BoomberLoad/cmd/.env .
EXPOSE 8080
ENTRYPOINT [ "./boomberload" ]