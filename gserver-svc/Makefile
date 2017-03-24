SOURCEDIR := .
SOURCES := $(shell find . -maxdepth 1 -type f \( -iname "*.go" ! -iname "*_test.*" \))
BIN_DIR := bin
BINARY := gserver-svc
FLAGS := -a -ldflags "-w"
OS := linux
DOCKERFILE := Dockerfile
DOCKER_IMAGE_NAME := gserver-svc

all: build docker-build 

build: $(SOURCES)
	if [ ! -d ${BIN_DIR} ]; then mkdir -p ${BIN_DIR} ; fi 
	CGO_ENABLED=0 GOOS=${OS} go build -o ${BIN_DIR}/${BINARY} ${FLAGS} ${SOURCES}


docker-build: $(DOCKERFILE)
	docker build -f ${DOCKERFILE} -t ${DOCKER_IMAGE_NAME} ${SOURCEDIR} 

.PHONY:
clean:
	if [ -f ${BIN_DIR}/${BINARY} ] ; then rm ${BIN_DIR}/${BINARY} ; fi