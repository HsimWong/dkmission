IMAGE_VERSION = latest
REGISTRY = docker.io/hsimwong
IMAGE = ${REGISTRY}/dkmission:${IMAGE_VERSION}

#.PHONY: build
all: manager worker

manager:
	GOOS=linux go build   -o build/manager app/manager_main.go

worker:
	GOOS=linux go build   -o build/worker app/worker_main.go

dockerize:
	docker build -t ${IMAGE} .


pushImage:
	docker push ${IMAGE}

clean:
	rm -rf ./build
