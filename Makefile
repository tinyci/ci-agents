VERSION=$(shell cat VERSION)
CONTAINER_DIR=/go/src/github.com/tinyci/ci-agents

GOLANGCI_LINT_VERSION=1.39.0

TESTRUN ?=
TESTPATH ?= ./...

STD_DOCKERFILE=dockerfiles/Dockerfile
RELEASE_DOCKERFILE=dockerfiles/Dockerfile.release
RELEASE_CONTEXT=release

DOCKER_RUN=docker run \
					 --rm
DOCKER_CONTAINER_DIR=-v ${PWD}:$(CONTAINER_DIR) \
								-w $(CONTAINER_DIR)

DOCKER_BUILD_CACHE_VOLUME=-v ci-agents-buildcache:/go/pkg
DEMO_DOCKER_IMAGE=tinyci-agents
DEBUG_DOCKER_IMAGE=tinyci-agents-debug
TEST_DOCKER_IMAGE=tinyci-agents-test
BUILD_DOCKER_IMAGE=tinyci-build
GO_VERSION=1.16

DEBUG_PORTS= -p 3000:3000 \
								-p 6000:6000 \
								-p 6001:6001 \
								-p 6002:6002 \
								-p 6005:6005 \
								-p 6010:6010

BUILD_DOCKER_RUN=\
								$(DOCKER_RUN) \
								-e GOBIN=/tmp/bin/tinyci-$(VERSION) \
								-e GOCACHE=/tmp/cache \
								-u $$(id -u):$$(id -g) \
								-v ${PWD}/build/:/tmp/bin \
								$(DOCKER_BUILD_CACHE_VOLUME) \
								-w $(CONTAINER_DIR) \
								-v ${PWD}:$(CONTAINER_DIR) \
								golang:${GO_VERSION}

TEST_DOCKER_RUN=\
								$(DOCKER_RUN) -it \
								-e CREATE_DB=1 \
								-e GIN_MODE=test \
								-e TESTING=1 \
								-e TESTRUN='${TESTRUN}' \
								-e TESTPATH='${TESTPATH}' \
								$(DOCKER_BUILD_CACHE_VOLUME) \
								--name $(TEST_DOCKER_IMAGE) \
								$(DOCKER_CONTAINER_DIR) \
								$(TEST_DOCKER_IMAGE)

DEBUG_DOCKER_RUN=\
								$(DOCKER_RUN) -it \
								-e CREATE_DB=1 \
								-e DEBUG=1 \
								$(if ${USE_JAEGER}, -e JAEGER_AGENT_HOST=jaegertracing,) \
								$(DEBUG_PORTS) \
								--link react:react \
								$(DOCKER_BUILD_CACHE_VOLUME) \
								$(if ${USE_JAEGER}, --link jaegertracing:jaegertracing,) \
								--name $(DEBUG_DOCKER_IMAGE) \
								$(DOCKER_CONTAINER_DIR) \
								$(DEBUG_DOCKER_IMAGE)

DEMO_DOCKER_RUN=\
								$(DOCKER_RUN) -it \
								-v ${PWD}/.ca:/var/ca \
								-v ${PWD}/.db:/var/lib/postgresql \
								-v ${PWD}/.logs:/var/tinyci/logs \
								-e START_SERVICES="${START_SERVICES}" \
								$(if ${USE_JAEGER}, -e JAEGER_AGENT_HOST=jaegertracing,) \
								-e DEBUG=1 \
								$(DEBUG_PORTS) \
								--link react:react \
								$(if ${USE_JAEGER}, --link jaegertracing:jaegertracing,) \
								$(DOCKER_CONTAINER_DIR) \
								--name $(DEMO_DOCKER_IMAGE) \
								$(DEMO_DOCKER_IMAGE)

test: build-image
	$(TEST_DOCKER_RUN) make 'TESTRUN=${TESTRUN}' 'TESTPATH=${TESTPATH}' do-test

do-test:
	go test -v -timeout 30m -p 1 -race ${TESTPATH} -check.v -check.f "${TESTRUN}" # -p 1 is needed because of gorilla/sessions init routines

test-debug: build-image
	$(TEST_DOCKER_RUN) bash

test-debug-attach:
	docker exec -it $(DEBUG_DOCKER_IMAGE) bash

demo-shell:
	docker-compose exec tinyci bash

demo-sql-shell:
	docker exec -it $(DEMO_DOCKER_IMAGE) psql tinyci

do-build:
	GO111MODULE=on go install -v -ldflags "-X main.TinyCIVersion=$(VERSION)" ./cmd/... ./api/...
	cp .config/services.yaml.example $${GOBIN:-${GOPATH}/bin}

build: distclean
	mkdir -p build
	docker pull golang:${GO_VERSION}
	$(BUILD_DOCKER_RUN) make do-build

distclean:
	rm -rf build ${RELEASE_CONTEXT}

dist: build
	mkdir -p ${RELEASE_CONTEXT}
	tar -C build -cvzf ${RELEASE_CONTEXT}/tinyci-$(VERSION).tar.gz tinyci-$(VERSION)

release: distclean dist
	docker build --build-arg VERSION="${VERSION}" -t "tinyci/release:${VERSION}" -f ${RELEASE_DOCKERFILE} ${RELEASE_CONTEXT}

demo: stop-demo build-demo-image
	docker-compose up

stop-demo:
	docker-compose rm -f

clean-demo: build-demo-image stop-demo
	$(DOCKER_RUN) --entrypoint /bin/bash -v ${PWD}/.ca:/var/ca -v ${PWD}/.logs:/var/tinyci/logs -v ${PWD}/.db:/var/lib/postgresql $(DEMO_DOCKER_IMAGE) -c "rm -rf /var/lib/postgresql/11; rm -rf /var/tinyci/logs/*; rm -rf /var/ca/*"

build-demo-image:
	docker build -t $(DEMO_DOCKER_IMAGE) -f $(STD_DOCKERFILE) .

update-demo-image:
	docker build --no-cache -t $(DEBUG_DOCKER_IMAGE) -f $(STD_DOCKERFILE) .

build-image:
	docker build --build-arg TESTING=1 -t $(TEST_DOCKER_IMAGE) -f $(STD_DOCKERFILE) .

tag-test-image:
	docker build --build-arg TESTING=1 --no-cache -t tinyci/ci-agents:$(shell date '+%m.%d.%Y') -f $(STD_DOCKERFILE) .

update-task-ymls:
	sed -i -e 's!^default_image: tinyci/ci-agents:.*$$!default_image: tinyci/ci-agents:$(shell date '+%m.%d.%Y')!g' $$(find . -name task.yml)

update-modules:
	rm -rf go.mod go.sum
	GO111MODULE=on go get -u -d ./...
	GO111MODULE=on go mod tidy

check-service-config:
	if [ ! -f .config/services.yaml ]; \
	then \
	  echo \
	  echo 2>&1 "Please create .config/services.yaml from the example in .config/services.yaml.example" \
	  echo \
	  exit 1; \
	fi

start-selective-services:
	for srv in ${START_SERVICES}; do pkill -f $$srv; (tinyci service $$srv &); done

start-services: check-service-config
	pkill tinyci || :
	go install -v ./cmd/...
	@if [ "x${START_SERVICES}" != "x" ]; then make start-selective-services; exit 0; fi
	tinyci launch

wait:
	sleep infinity

bin/golangci-lint:
	mkdir -p bin
	wget -O- https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz | tar vxz --strip-components=1 -C bin golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64/golangci-lint

golangci-lint: bin/golangci-lint
	bin/golangci-lint run -v

gen: mockgen build-image
	cd ci-gen && make gen	
	go generate -v ./db/migrations
	$(TEST_DOCKER_RUN) bash -c "go generate ./..."


mockgen:
	GO111MODULE=off go get github.com/golang/mock/...
	${GOPATH}/bin/mockgen -package github github.com/tinyci/ci-agents/clients/github Client > mocks/github/mock.go

jaeger:
	docker run --name jaegertracing -idt -p 16686:16686 jaegertracing/all-in-one:latest --log-level debug || :

stop-jaeger:
	docker rm -f jaegertracing
