ROOT_PATH=/home/mobius/projects/ankr/dccn-tools/examples/
GOOGLE_PATH="/home/mobius/projects/google"
.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cpmain .
	docker build --no-cache -t registry.cn-hangzhou.aliyuncs.com/mobius/cpmain:0.1 .
	rm -rf cpmain
.PHONY: docker
docker:
	docker push registry.cn-hangzhou.aliyuncs.com/mobius/cpmain:0.1

.PHONY: proto
proto:
	protoc --proto_path=$(ROOT_PATH):. --ankr_out=plugins=ankr:. proto/checkpassword.proto