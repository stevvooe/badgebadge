FROM golang:1.8 as builder
WORKDIR /go/src/badge
COPY . /go/src/badge
RUN apt-get update; apt-get install git
RUN go get github.com/LK4D4/vndr
RUN vndr
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY templates /templates
COPY assets /assets
COPY badges /badges
 COPY . /go/src/app
COPY --from=builder /go/src/badge/app .
CMD ["./app"]

EXPOSE 8080


