# ------------------------------------ Lint ------------------------------------
FROM golangci/golangci-lint:v1.61-alpine AS lint 

ARG app_name

# Buildkit will skip stages the final image doesn't depend on.
# Creating and copying a dummy file to avoid buildkit skipping this stage.
RUN echo "Dummy" > /tmp/lint-dummy.txt

RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /app

COPY Taskfile.yaml .
COPY config.env* .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN task go:lint

# -------------------------------- License Check -------------------------------
FROM golang:1.23-alpine AS license

ARG app_name

RUN echo "Dummy" > /tmp/license-dummy.txt

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN go install github.com/google/go-licenses@latest

WORKDIR /app

COPY Taskfile.yaml .
COPY config.env* .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN task go:license-check

# ------------------------------------ Build -----------------------------------
FROM golang:1.23-alpine AS build

ARG app_name

RUN go install github.com/go-task/task/v3/cmd/task@latest

WORKDIR /app

COPY Taskfile.yaml .
COPY config.env* .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN task go:test
RUN task go:build

# ----------------------------------- Final -----------------------------------
# Using this instead of scratch with CA certificates and tzdata. This is slightly
# bigger (about 1MB) but makes for a cleaner setup.
FROM gcr.io/distroless/static-debian12:nonroot AS final

ARG artifacts

# Creating and copying a dummy file to avoid buildkit skipping this stage.
# The reason for copying the file here, in stead of the build stage, is that the
# build and lint stage can then be run in parallell.
COPY --from=lint /tmp/lint-dummy.txt /tmp/lint-dummy.txt 
COPY --from=license /tmp/license-dummy.txt /tmp/license-dummy.txt 

COPY --from=build --chown=nonroot:nonroot /app/${artifacts}/app.bin /app

ENTRYPOINT ["/app"]
