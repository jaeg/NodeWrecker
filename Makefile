REPO = jaeg/nodewrecker
VERSION = 0.0.4

image: build-linux
	docker build -t $(REPO):$(VERSION) .

image-pi: build-linux-pi
	docker build -t $(REPO):$(VERSION)-pi .

build:
	go build -o pkg/NodeWrecker

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o pkg/NodeWrecker-linux

build-linux-pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o pkg/NodeWrecker-linux-pi

publish-pi:
	docker push $(REPO):$(VERSION)-pi