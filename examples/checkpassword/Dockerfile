FROM registry.cn-hangzhou.aliyuncs.com/mobius/alpine:v1.0.3 AS builder
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY ./ ./
RUN rm -f Dockerfile
RUN protoc --ankr_out=plugins=ankr:. proto/checkpassword.proto
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/app .
CMD ["/root/app"]