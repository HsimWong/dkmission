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
	rm processor.o
	make clean -C darknet 
	sudo rm $(PREFIX)/lib/libdarknet.so
	sudo rm $(PREFIX)/include/darknet.h

dependency: darknetinstall
	sudo apt install imagemagick -y
	cd processor && gcc -c processor.c -ldarknet

darknet-prepare:
	git submodule sync --recursive
	make -C darknet -j8

darknetinstall: darknet-prepare 
	sudo cp darknet/libdarknet.so $(PREFIX)/lib
	sudo cp darknet/include/darknet.h $(PREFIX)/include
	sudo ldconfig
