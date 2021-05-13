# detect operating system
ifeq ($(OS),Windows_NT)
    CURRENT_OS := Windows
else
    CURRENT_OS := $(shell uname -s)
endif

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -tags netgo -ldflags "-X github.com/deep2chain/sscq/version.GitCommit=${COMMIT_HASH}"
# docker part[single-node part]
DOCKER_TAG = unstable
DOCKER_IMAGE = falcon0125/sscqsingle

build: clean
ifeq ($(CURRENT_OS),Windows)
	go build  $(BUILD_FLAGS) -o ../../build/ssd.exe ../../cmd/ssd
	go build  $(BUILD_FLAGS) -o ../../build/sscli.exe ../../cmd/sscli
else
	go build  $(BUILD_FLAGS) -o ../../build/ssd ../../cmd/ssd
	go build  $(BUILD_FLAGS) -o ../../build/sscli ../../cmd/sscli
endif

clean:
	@rm -rf ../../build

docker: build
	@docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	@docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	@docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_HASH}

run:
	@docker run --rm -it ${DOCKER_IMAGE}:latest ssd start --minimum-gas-prices=0.0001sscq

shell:
	@docker run --rm -it ${DOCKER_IMAGE}:latest sh
