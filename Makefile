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
								-v ${PWD}/build/tinyci-$(VERSION):/build \
								-w $(CONTAINER_DIR) \
								-e GOBIN=/build \
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

demo-sql-shell:
	docker exec -it $(DEMO_DOCKER_IMAGE) psql tinyci

do-build:
	go get github.com/erikh/migrator
	go install -v -ldflags "-X main.TinyCIVersion=$(VERSION)" ./...
	cp .config/services.yaml.example $${GOBIN:-${GOPATH}/bin}
	cp -Rfp migrations $${GOBIN:-${GOPATH}/bin}

build: build-build-image
	$(BUILD_DOCKER_RUN) make do-build

dist: build
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

gen: build-demo-image
	$(DOCKER_RUN) -u $$(id -u):$$(id -g) $(DOCKER_CONTAINER_DIR) --entrypoint /bin/bash $(DEMO_DOCKER_IMAGE) gen.sh

gen-javascript:
	mkdir -p $${TARGET_DIR:-${PWD}/gen/javascript}
	docker run --rm -u $$(id -u):$$(id -g) -v $${TARGET_DIR:-${PWD}/gen/javascript}:/swagger -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate \
		-i /local/swagger/uisvc/swagger.yml \
		-l javascript \
		-o /swagger

swagger-serve:
	docker run -p 8080:8080 -it -v ${PWD}:/swagger redoc-cli serve file:///swagger/uisvc/swagger.yml

swagger-docs:
	docker run -it -u $(shell id -u):$(shell id -g) -v ${PWD}/swagger:/swagger redoc-cli bundle file:///swagger/uisvc/swagger.yml -o /swagger/docs.html

check-s3cmd:
	@which s3cmd 2>&1 >/dev/null || echo "You must install a working copy of s3cmd configured to upload to the docs.tinyci.org bucket."

grpc-docs: build-debug-image
	mkdir -p grpc/docs
	$(DEBUG_DOCKER_RUN) bash -c "protoc --doc_out=grpc/docs --doc_opt=html,index.html --proto_path=/go/src $(CONTAINER_DIR)/grpc/services/**/*.proto $(CONTAINER_DIR)/grpc/types/*.proto"

upload-docs: check-s3cmd swagger-docs grpc-docs
	s3cmd put swagger/docs.html -m text/html s3://docs.tinyci.org/swagger/index.html
	s3cmd put grpc/docs/index.html -m text/html s3://docs.tinyci.org/grpc/index.html

swagger-validate: require-spec build-demo-image
	$(DOCKER_RUN) -u $$(id -u):$$(id -g) $(DOCKER_CONTAINER_DIR) --entrypoint /go/bin/swagger $(DEMO_DOCKER_IMAGE) \
		validate swagger/$${SPEC}/swagger.yml

staticcheck:
	go get honnef.co/go/tools/...
	staticcheck ./...
