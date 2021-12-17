Branch=$(shell git symbolic-ref --short -q HEAD)
Commit=$(shell git rev-parse --short HEAD)
Date=$(shell git log --pretty=format:%cd $(Commit) -1)
Author=$(shell git log --pretty=format:%an $(Commit) -1)
shortDate=$(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%Y%m%d)
Email=$(shell git log --pretty=format:%ae $(Commit) -1)
Ver=$(shell echo $(Branch)-$(Commit)-$(shortDate))
GoVersion=$(shell go version )
Repo ?= ankrnetwork/nsq-clear

.PHONY: build
build: fmt vet
	CGO_ENABLED=0 go build -o build/nsq-clear -a -installsuffix cgo \
    -ldflags "-X 'github.com/Ankr-network/aci-syncer/cmd.Branch=$(Branch)' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Commit=$(Commit)' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Date=$(Date)' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Author=$(Author)' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.Email=$(Email)' \
    -X 'github.com/Ankr-network/aci-syncer/cmd.GoVersion=$(GoVersion)'"


.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vat
vet:
	go vet ./...

.PHONY: docker
docker:
	docker build --build-arg Branch='$(Branch)' \
		--build-arg Commit='$(Commit)' \
		--build-arg Date='$(Date)' \
		--build-arg Author='$(Author)' \
		--build-arg Email='$(Email)' \
		-t $(Repo):$(Commit) .
	docker image prune -f --filter label=stage=builder

.PHONY: release
release: docker
	docker push $(Repo):$(Commit)

.PHONY: run
run:
	./build/nsq-clear \
	--nsq_admin "http://127.0.0.1:4171" \
	--threshold 100 \
	--schedule 1m \
	--level "debug"