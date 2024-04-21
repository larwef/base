# TODO: Uncomment if using a seperate file for envronemtn variables you don't
# want to disclose.
# include config.env

# TODO: Change name. Remember to also change the name of cmd/app/.
APP_NAME=app
VERSION=v0.0.1

# TDOD: Change registry.
REGISTRY=your.registry/name

# TODO: Change OS and ARCH.
GOOS=linux
GOARCH=amd64

ARTIFACTS=./artifacts

clean:
	rm -rf $(ARTIFACTS)

# ------------------------------------- Go -------------------------------------
run:
	make generate
	LOG_LEVEL=debug \
	LOG_JSON=false \
		go run cmd/app/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	GOOS=$(GOOS) \
    GOARCH=$(GOARCH) \
    go build \
    	-ldflags " \
			-X main.appName=$(APP_NAME) \
    		-X main.version=$(VERSION) \
    	" -o $(ARTIFACTS)/app.bin cmd/$(APP_NAME)/main.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

# ----------------------------------- Docker -----------------------------------
docker-build:
	docker build -t $(REGISTRY)/$(APP_NAME):$(VERSION) \
		--build-arg app_name=$(APP_NAME) \
		--build-arg artifacts=$(ARTIFACTS) \
		-f build/package/Dockerfile .

docker-run:
	docker run -it --rm \
		--name $(APP_NAME) \
		--env=ADDRESS=:8080 \
		-p 8080:8080 \
		$(REGISTRY)/$(APP_NAME):$(VERSION)

docker-push:
	docker push $(REGISTRY)/$(APP_NAME):$(VERSION)
