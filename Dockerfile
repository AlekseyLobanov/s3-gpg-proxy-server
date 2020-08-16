FROM golang:1.14 AS builder
WORKDIR /root/
RUN GO111MODULE=on go get github.com/minio/minio-go/v7  
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /root/app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates gnupg
WORKDIR /root/
COPY --from=builder /root/app /root/app
CMD ["/root/app"]