SOURCEDIR := .
SOURCES := $(shell find . -maxdepth 1 -type f \( -iname "*.go" ! -iname "*_test.*" \))
BIN_DIR := bin
PROTO_DIR := auth
BINARY := auth-svc
FLAGS := -a -ldflags "-w"
OS := linux
DOCKERFILE := Dockerfile
DOCKER_IMAGE_NAME := auth-svc

all: build docker-build 

build: $(SOURCES)
	if [ ! -d ${BIN_DIR} ]; then mkdir -p ${BIN_DIR} ; fi 
	protoc --go_out=plugins=grpc:${SOURCEDIR} ${PROTO_DIR}/*.proto && \
	CGO_ENABLED=0 GOOS=${OS} go build -o ${BIN_DIR}/${BINARY} ${FLAGS} ${SOURCES}

docker-build: $(DOCKERFILE)
	docker build -f ${DOCKERFILE} -t ${DOCKER_IMAGE_NAME} ${SOURCEDIR} 

.PHONY:
clean:
	if [ -d ${BIN_DIR} ] ; then rm -r ${BIN_DIR} ; fi
	if [ -f ${PROTO_DIR}/*.pb.go ] ; then rm ${PROTO_DIR}/*.pb.go ; fi
