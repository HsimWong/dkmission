IMAGE_VERSION = latest
REGISTRY = docker.io/hsimwong
IMAGE = ${REGISTRY}/dkmission:${IMAGE_VERSION}
PREFIX = /usr/local


#.PHONY: build
all: manager worker

manager: dependency
	GOOS=linux LDFLAGS=-ldarknet go build   -o build/manager app/manager_main.go

worker: dependency
	GOOS=linux LDFLAGS=-ldarknet go build   -o build/worker app/worker_main.go 

dockerize:
	docker build -t ${IMAGE} .


pushImage:
	docker push ${IMAGE}

clean:
	rm -rf ./build
	make clean -C darknet
	rm processor/processor.o 
	rm processor/libdarknet.so
	rm processor/darknet.h
	rm ${PREFIX}/include/darknet.h

dependency: darknetinstall
	apt install imagemagick -y
	cd processor && gcc -c processor.c -L. -ldarknet -Wl,-rpath,$PWD/libdarknet.so

darknet-prepare:
	git submodule update --init --recursive
	make -C darknet -j8

darknetinstall: darknet-prepare
	cp darknet/libdarknet.so processor
	cp darknet/include/darknet.h processor
	cp darknet/include/darknet.h /usr/local/include/

dockerPrepare:
	cp darknet/libdarknet.so processor
	cp darknet/include/darknet.h processor
	cp darknet/include/darknet.h /usr/local/include/
