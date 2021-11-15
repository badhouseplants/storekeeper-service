FROM golang:1.17.3-alpine3.14 as builder
WORKDIR /go/src/app
COPY . /go/src/app
RUN apk add git
RUN go build -i main.go


FROM  alpine:3.14.0
WORKDIR /root/
COPY migrations/scripts/ /root/migrations/scripts/ 
COPY --from=builder /go/src/app/main .
CMD ["./main"]  