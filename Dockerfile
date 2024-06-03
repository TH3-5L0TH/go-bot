FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go build -o /go/bin/app .


FROM alpine:latest

COPY --from=builder /go/bin/app /usr/local/bin/app

CMD ["app"]
