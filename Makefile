# TODO: Change name. Remember to also change the name of cmd/app/.
APP_NAME=app
VERSION=v0.0.1

# TDOD: Change registry.
REGISTRY=your.registry/name

# TODO: Change OS and ARCH.
GOOS=linux
GOARCH=amd64

# ------------------------------------- Go -------------------------------------
run:
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
    	" -o app.bin cmd/$(APP_NAME)/main.go

# ----------------------------------- Docker -----------------------------------
docker-build:
	docker build -t $(REGISTRY)/$(APP_NAME):$(VERSION) \
		--build-arg app_name=$(APP_NAME) \
		-f build/package/Dockerfile .

docker-run:
	docker run -it --rm \
		--name $(APP_NAME) \
		--env=ADDRESS=:8080 \
		-p 8080:8080 \
		$(REGISTRY)/$(APP_NAME):$(VERSION)

docker-push:
	docker push $(REGISTRY)/$(APP_NAME):$(VERSION)
