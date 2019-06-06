VERSION=$(shell cat VERSION)
CONTAINER_DIR=/go/src/github.com/tinyci/ci-agents/ci-gen
DOCKER_RUN=docker run \
					 --rm -it
DOCKER_CONTAINER_DIR=-v ${PWD}:$(CONTAINER_DIR) \
								-w $(CONTAINER_DIR)

DOCKER_IMAGE=tinyci-gen
DOCKER_RUN_FULL=\
								$(DOCKER_RUN) \
								-v ${PWD}/build/:/build \
								-w $(CONTAINER_DIR) \
								-e GOBIN=/build/tinyci-$(VERSION) \
								$(DOCKER_IMAGE)

SWAGGER_SERVICES := uisvc

get-box:
	@if [ ! -f "$(shell which box)" ]; \
	then \
		echo "Need to install box to build the docker images we use. Requires root access."; \
		curl -sSL box-builder.sh | sudo bash; \
	fi

build-image: get-box
	box -t $(DOCKER_IMAGE) box.rb

gen: build-image
	$(DOCKER_RUN) -u $$(id -u):$$(id -g) $(DOCKER_CONTAINER_DIR) $(DOCKER_IMAGE) bash gen.sh

gen-javascript:
	mkdir -p $${TARGET_DIR:-${PWD}/gen/javascript}
	docker run --rm -u $$(id -u):$$(id -g) -v $${TARGET_DIR:-${PWD}/gen/javascript}:/swagger -v ${PWD}:/local swaggerapi/swagger-codegen-cli generate \
		-i /local/swagger/uisvc/swagger.yml \
		-l javascript \
		-o /swagger

swagger-serve:
	docker run -p 8080:8080 -it -v ${PWD}:/swagger tinyci/redoc-cli serve file:///swagger/uisvc/swagger.yml

swagger-docs:
	docker run --rm -it -u $(shell id -u):$(shell id -g) -v ${PWD}/swagger:/swagger tinyci/redoc-cli bundle file:///swagger/uisvc/swagger.yml -o /swagger/docs.html

grpc-docs: build-image
	mkdir -p grpc/docs
	$(DOCKER_RUN) --rm $(DOCKER_CONTAINER_DIR) --entrypoint '' $(DOCKER_IMAGE) bash -c "protoc --doc_out=grpc/docs --doc_opt=html,index.html --proto_path=/go/src $(CONTAINER_DIR)/grpc/services/**/*.proto $(CONTAINER_DIR)/grpc/types/*.proto"
