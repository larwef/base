# TODO: Change name. Remember to also change the name of cmd/app/.
APP_NAME=app
VERSION=v0.0.1

# TDOD: Change registry.
REGISTRY=your.registry/name

# TODO: Change OS and ARCH.
GOOS=linux
GOARCH=amd64

# TODO: Change to match your Kubernetes setup.
K8S_CONTEXT=your_context
K8S_NAMESPACE=your_namespace

ARTIFACTS=artifacts

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
    	" -o $(ARTIFACTS)/app.bin cmd/$(APP_NAME)/main.go

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

# --------------------------------- Kubernetes ---------------------------------
# Remmeber to add more variables here if you add more to be susbtituted in the
# manifest templates.
k8s-plan:
	APP_NAME=$(APP_NAME) \
	VERSION=$(VERSION) \
	REGISTRY=$(REGISTRY) \
	K8S_CONTEXT=$(K8S_CONTEXT) \
	K8S_NAMESPACE=$(K8S_NAMESPACE) \
	MANIFEST_OUTPUT=$(ARTIFACTS)/k8s-manifest.yaml \
		./scripts/k8s-plan.sh

k8s-apply:
	K8S_CONTEXT=$(K8S_CONTEXT) \
	K8S_NAMESPACE=$(K8S_NAMESPACE) \
		./scripts/k8s-apply.sh

# Included for completeness. Commented for safety.
# kubernetes-delete:
# 	K8S_CONTEXT=$(K8S_CONTEXT) \
# 	K8S_NAMESPACE=$(K8S_NAMESPACE) \
# 		./scripts/delete.sh