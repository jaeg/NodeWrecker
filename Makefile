REPO = jaeg/nodewrecker
VERSION = 1.0.0

image: build-linux
	docker build -t $(REPO):$(VERSION) . --build-arg binary=NodeWrecker-linux --build-arg version=$(VERSION)

image-pi: build-linux-pi

	docker build -t $(REPO):$(VERSION)-pi . --build-arg binary=NodeWrecker-linux-pi --build-arg version=$(VERSION)

build:
	go build -o pkg/NodeWrecker

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o pkg/NodeWrecker-linux

build-linux-pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o pkg/NodeWrecker-linux-pi

publish-pi:
	docker push $(REPO):$(VERSION)-pi
	docker tag $(REPO):$(VERSION)-pi $(REPO):latest-pi
	docker push $(REPO):latest-pi

publish:
	docker push $(REPO):$(VERSION)
	docker tag $(REPO):$(VERSION) $(REPO):latest
	docker push $(REPO):latest