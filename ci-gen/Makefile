VERSION=$(shell cat VERSION)
CONTAINER_DIR=/go/src/github.com/tinyci/ci-agents
DOCKER_RUN=docker run \
					 --rm -it
DOCKER_CONTAINER_DIR=-v $(shell realpath ${PWD}/..):$(CONTAINER_DIR) \
								-w $(CONTAINER_DIR)
DOCKER_IMAGE=tinyci-gen
SWAGGER_SERVICES := uisvc
TARGET_DIR=${PWD}/ci-gen

get-box:
	@if [ ! -f "$(shell which box)" ]; \
	then \
		echo "Need to install box to build the docker images we use. Requires root access."; \
		curl -sSL box-builder.sh | sudo bash; \
	fi

build-gen-image: get-box
	box -t $(DOCKER_IMAGE) box.rb

gen: build-gen-image
	$(DOCKER_RUN) -u $$(id -u):$$(id -g) $(DOCKER_CONTAINER_DIR) $(DOCKER_IMAGE) bash ci-gen/gen.sh

gen-javascript:
	rm -rf ${TARGET_DIR}/gen/javascript
	mkdir -p ${TARGET_DIR}/gen/javascript
	docker pull openapitools/openapi-generator-cli:latest
	docker run --rm -u $$(id -u):$$(id -g) -v ${TARGET_DIR}/gen/javascript:/swagger -v ${PWD}:/local openapitools/openapi-generator-cli generate \
		-i /local/swagger/uisvc/swagger.yml \
		-g javascript \
		-o /swagger

swagger-serve:
	docker run -p 8080:8080 -it -v ${PWD}:/swagger tinyci/redoc-cli serve file:///swagger/uisvc/swagger.yml

swagger-docs:
	docker run --rm -it -u $(shell id -u):$(shell id -g) -v ${PWD}/swagger:/swagger tinyci/redoc-cli bundle file:///swagger/uisvc/swagger.yml -o /swagger/docs.html

grpc-docs: build-gen-image
	mkdir -p ci-gen/grpc/docs
	$(DOCKER_RUN) --rm $(DOCKER_CONTAINER_DIR) --entrypoint '' $(DOCKER_IMAGE) bash -c "protoc --doc_out=ci-gen/grpc/docs --doc_opt=html,index.html --proto_path=/go/src $(CONTAINER_DIR)/grpc/services/**/*.proto $(CONTAINER_DIR)/grpc/types/*.proto"
