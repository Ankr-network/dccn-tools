FROM  golang:1.16.7-alpine3.14 AS builder
LABEL stage=builder
RUN mkdir /go/src/app
RUN go env -w GOPROXY="https://goproxy.cn,direct"
WORKDIR /go/src/app
COPY ./ ./
ARG Branch=""
ARG Commit=""
ARG Date=""
ARG Author=""
ARG Email=""
RUN export GoVersion=$(go version) && CGO_ENABLED=0 GOOS=linux go build -o clear -a -installsuffix cgo -ldflags \
    "-X 'github.com/Ankr-network/aci-syncer/cmd.Branch=${Branch}' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Commit=${Commit}' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Date=${Date}' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Author=${Author}' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Email=${Email}' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.GoVersion=${GoVersion}'" .

FROM alpine:3.14
WORKDIR /root/
COPY --from=builder /go/src/app/clear .
ENTRYPOINT ["/root/clear"]
