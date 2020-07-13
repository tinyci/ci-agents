VERSION=$(shell cat VERSION)
CONTAINER_DIR=/go/src/github.com/tinyci/ci-agents

STD_BOXFILE=box-builds/box.rb
RELEASE_BOXFILE=box-builds/box-release.rb

DOCKER_RUN=docker run \
					 --rm
DOCKER_CONTAINER_DIR=-v ${PWD}:$(CONTAINER_DIR) \
								-w $(CONTAINER_DIR)

DEMO_DOCKER_IMAGE=tinyci-agents
DEBUG_DOCKER_IMAGE=tinyci-agents-debug
TEST_DOCKER_IMAGE=tinyci-agents-test
BUILD_DOCKER_IMAGE=tinyci-build
GO_VERSION=1.14

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
								-w $(CONTAINER_DIR) \
								-v ${PWD}:$(CONTAINER_DIR) \
								golang:${GO_VERSION}

TEST_DOCKER_RUN=\
								$(DOCKER_RUN) -it \
								-e CREATE_DB=1 \
								-e GIN_MODE=test \
								-e TESTING=1 \
								--name $(TEST_DOCKER_IMAGE) \
								$(DOCKER_CONTAINER_DIR) \
								$(TEST_DOCKER_IMAGE)

DEBUG_DOCKER_RUN=\
								$(DOCKER_RUN) -it \
								-e CREATE_DB=1 \
								-e DEBUG=1 \
								-e JAEGER_AGENT_HOST=jaegertracing \
								$(DEBUG_PORTS) \
								--link react:react \
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
								-e JAEGER_AGENT_HOST=jaegertracing \
								-e DEBUG=1 \
								$(DEBUG_PORTS) \
								--link react:react \
								$(if ${USE_JAEGER}, --link jaegertracing:jaegertracing,) \
								$(DOCKER_CONTAINER_DIR) \
								--name $(DEMO_DOCKER_IMAGE) \
								$(DEMO_DOCKER_IMAGE)

SWAGGER_SERVICES := uisvc

test: build-image
	$(TEST_DOCKER_RUN) make do-test

do-test:
	go test -timeout 30m -p 1 -race -v ./... -check.v # -p 1 is needed because of gorilla/sessions init routines

test-debug: build-debug-image
	$(DEBUG_DOCKER_RUN) bash

test-debug-attach:
	docker exec -it $(DEBUG_DOCKER_IMAGE) bash

demo-shell:
	docker-compose exec tinyci bash

demo-sql-shell:
	docker exec -it $(DEMO_DOCKER_IMAGE) psql tinyci

do-build:
	GOPATH=$$(mktemp -d /tmp/gopath.XXXXX) go install -v github.com/erikh/migrator
	GO111MODULE=on go install -v -ldflags "-X main.TinyCIVersion=$(VERSION)" ./cmd/... ./api/...
	cp .config/services.yaml.example $${GOBIN:-${GOPATH}/bin}
	cp -Rfp migrations $${GOBIN:-${GOPATH}/bin}

build: distclean
	mkdir -p build
	docker pull golang:${GO_VERSION}
	$(BUILD_DOCKER_RUN) make do-build

distclean:
	rm -rf build tinyci-$(VERSION).tar.gz

dist: build
	tar -C build -cvzf tinyci-$(VERSION).tar.gz tinyci-$(VERSION)

release: distclean dist
	VERSION="$(VERSION)" box -t "tinyci/release:$(VERSION)" $(RELEASE_BOXFILE)

demo: stop-demo build-demo-image
	docker-compose up

stop-demo:
	docker-compose rm -f

clean-demo: build-demo-image stop-demo
	$(DOCKER_RUN) --entrypoint /bin/bash -v ${PWD}/.ca:/var/ca -v ${PWD}/.logs:/var/tinyci/logs -v ${PWD}/.db:/var/lib/postgresql $(DEMO_DOCKER_IMAGE) -c "rm -rf /var/lib/postgresql/11; rm -rf /var/tinyci/logs/*; rm -rf /var/ca/*"

build-demo-image: get-box
	box -t $(DEMO_DOCKER_IMAGE) $(STD_BOXFILE)

build-debug-image: get-box
	DEBUG=1 box -t $(DEBUG_DOCKER_IMAGE) $(STD_BOXFILE)

update-demo-image: get-box
	DEBUG=1 box -n -t $(DEBUG_DOCKER_IMAGE) $(STD_BOXFILE)

update-image: get-box
	TESTING=1 box -t $(TEST_DOCKER_IMAGE) -n $(STD_BOXFILE)

build-image: get-box
	TESTING=1 box -t $(TEST_DOCKER_IMAGE) $(STD_BOXFILE)

tag-test-image: get-box
	PACKAGE_FOR_CI=1 TESTING=1 box -n -t tinyci/ci-agents:$(shell date '+%m.%d.%Y') $(STD_BOXFILE)

update-task-ymls:
	sed -i -e 's!^default_image: tinyci/ci-agents:.*$$!default_image: tinyci/ci-agents:$(shell date '+%m.%d.%Y')!g' $$(find . -name task.yml)

get-box:
	@if [ ! -f "$(shell which box)" ]; \
	then \
		echo "Need to install box to build the docker images we use. Requires root access."; \
		curl -sSL box-builder.sh | sudo bash; \
	fi

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
	for srv in ${START_SERVICES}; do pkill $$srv || :; ($$srv &); done

start-services: check-service-config
	pkill uisvc-server || :
	pkill logsvc || :
	pkill hooksvc || :
	pkill assetsvc || :
	pkill queuesvc || :
	pkill github-authsvc || :
	pkill datasvc || :
	go install -v ./cmd/... ./api/...
	@if [ "x${START_SERVICES}" != "x" ]; then make start-selective-services; exit 0; fi
	logsvc &
	assetsvc &
	datasvc &
	github-authsvc &
	queuesvc &
	uisvc-server &
	hooksvc &
	make wait

wait:
	sleep infinity

golangci-lint:
	go get github.com/golangci/golangci-lint/...
	golangci-lint run

gen: mockgen
	cd ci-gen && make gen	

mockgen:
	GO111MODULE=off go get github.com/golang/mock/...
	${GOPATH}/bin/mockgen -package github github.com/tinyci/ci-agents/clients/github Client > mocks/github/mock.go

jaeger:
	docker run --name jaegertracing -idt -p 16686:16686 jaegertracing/all-in-one:latest --log-level debug || :

stop-jaeger:
	docker rm -f jaegertracing
