IMAGE_VERSION = latest
REGISTRY = docker.io/hsimwong
IMAGE = ${REGISTRY}/dkmission:${IMAGE_VERSION}
PREFIX = /usr/local


#.PHONY: build
all: manager worker

manager: dependency
	GOOS=linux go build   -o build/manager app/manager_main.go

worker: dependency
	GOOS=linux go build   -o build/worker app/worker_main.go

dockerize:
	docker build -t ${IMAGE} .


pushImage:
	docker push ${IMAGE}

clean:
	rm -rf ./build
	rm $(PREFIX)/lib/libdarknet.so
	rm $(PREFIX)/include/darknet.h

dependency: darknetinstall
	sudo apt install imagemagick -y
	cd processor && gcc -c processor.c -L. -ldarknet -Wl,-rpath,$PWD/libdarknet.so

darknetinstall:
	cp processor/libdarknet.so $(PREFIX)/lib
	cp processor/darknet.h $(PREFIX)/include
	sudo ldconfig
