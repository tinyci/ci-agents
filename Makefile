VERSION=$(shell cat VERSION)
CONTAINER_DIR=/go/src/github.com/tinyci/ci-agents
DOCKER_RUN=docker run \
					 --rm -it
DOCKER_CONTAINER_DIR=-v ${PWD}:$(CONTAINER_DIR) \
								-w $(CONTAINER_DIR)
DEMO_DOCKER_IMAGE=tinyci-agents
DEBUG_DOCKER_IMAGE=tinyci-agents-debug
TEST_DOCKER_IMAGE=tinyci-agents-test
BUILD_DOCKER_IMAGE=tinyci-build

DEBUG_PORTS= -p 3000:3000 \
								-p 6000:6000 \
								-p 6001:6001 \
								-p 6002:6002 \
								-p 6005:6005 \
								-p 6010:6010

BUILD_DOCKER_RUN=\
								$(DOCKER_RUN) \
								-v ${PWD}/build/:/build \
								-w $(CONTAINER_DIR) \
								-e GOBIN=/build/tinyci-$(VERSION) \
								$(BUILD_DOCKER_IMAGE)

TEST_DOCKER_RUN=\
								$(DOCKER_RUN) \
								-e CREATE_DB=1 \
								-e GIN_MODE=test \
								-e TESTING=1 \
								--name $(TEST_DOCKER_IMAGE) \
								$(DOCKER_CONTAINER_DIR) \
								$(TEST_DOCKER_IMAGE)

DEBUG_DOCKER_RUN=\
								$(DOCKER_RUN) \
								-e CREATE_DB=1 \
								-e DEBUG=1 \
								$(DEBUG_PORTS) \
								--link react:react \
								--name $(DEBUG_DOCKER_IMAGE) \
								$(DOCKER_CONTAINER_DIR) \
								$(DEBUG_DOCKER_IMAGE)

DEMO_DOCKER_RUN=\
								$(DOCKER_RUN) \
								-v ${PWD}/.ca:/var/ca \
								-v ${PWD}/.db:/var/lib/postgresql \
								-v ${PWD}/.logs:/var/tinyci/logs \
								-e START_SERVICES="${START_SERVICES}" \
								-e DEBUG=1 \
								$(DEBUG_PORTS) \
								--link react:react \
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
	docker exec -it $(DEMO_DOCKER_IMAGE) bash

demo-sql-shell:
	docker exec -it $(DEMO_DOCKER_IMAGE) psql tinyci

do-build:
	go get github.com/erikh/migrator
	go install -v -ldflags "-X main.TinyCIVersion=$(VERSION)" ./...
	cp .config/services.yaml.example $${GOBIN:-${GOPATH}/bin}
	cp -Rfp migrations $${GOBIN:-${GOPATH}/bin}

build: build-build-image
	$(BUILD_DOCKER_RUN) make do-build

distclean:
	$(BUILD_DOCKER_RUN) bash -c 'rm -rf /build/*'

dist: build-build-image distclean build
	tar -C build -cvzf tinyci-$(VERSION).tar.gz tinyci-$(VERSION)

demo: build-demo-image
	$(DEMO_DOCKER_RUN) make start-services

clean-demo: build-demo-image
	$(DOCKER_RUN) --entrypoint /bin/bash -v ${PWD}/.ca:/var/ca -v ${PWD}/.logs:/var/tinyci/logs -v ${PWD}/.db:/var/lib/postgresql $(DEMO_DOCKER_IMAGE) -c "rm -rf /var/lib/postgresql/9.6; rm -rf /var/tinyci/logs/*; rm -rf /var/ca/*"

build-build-image: get-box
	box -t $(BUILD_DOCKER_IMAGE) box-build.rb

build-demo-image: get-box
	box -t $(DEMO_DOCKER_IMAGE) box.rb

build-debug-image: get-box
	DEBUG=1 box -t $(DEBUG_DOCKER_IMAGE) box.rb

update-demo-image: get-box
	DEBUG=1 box -n -t $(DEBUG_DOCKER_IMAGE) box.rb

update-image: get-box
	TESTING=1 box -t $(TEST_DOCKER_IMAGE) -n box.rb

build-image: get-box
	TESTING=1 box -t $(TEST_DOCKER_IMAGE) box.rb

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

vendor:
	GO111MODULE=on go mod vendor

.PHONY: vendor

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
	pkill datasvc || :
	go install -v ./cmd/... ./gen/...
	@if [ "x${START_SERVICES}" != "x" ]; then make start-selective-services; exit 0; fi
	logsvc &
	assetsvc &
	datasvc &
	queuesvc &
	uisvc-server &
	hooksvc &
	make wait

wait:
	sleep infinity

staticcheck:
	go get honnef.co/go/tools/...
	staticcheck ./...

mockgen:
	GO111MODULE=off go get github.com/golang/mock/...
	${GOPATH}/bin/mockgen -package github github.com/tinyci/ci-gen/clients/github Client > mocks/github/mock.go
