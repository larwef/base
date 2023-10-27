# ------------------------------------ Lint ------------------------------------
FROM golangci/golangci-lint:v1.55-alpine as lint

ARG app_name

# Buildkit will skip stages the final image doesn't depend on.
# Creating and copying a dummy file to avoid buildkit skipping this stage.
RUN echo "Dummy" > /tmp/dummy.txt

RUN apk add --no-cache make

WORKDIR /app

COPY Makefile .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN make test
RUN make build

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN make lint

# ------------------------------------ Build -----------------------------------
FROM golang:1.21-alpine as build

ARG app_name

RUN apk add --no-cache make

WORKDIR /app

COPY Makefile .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/${app_name}/main.go cmd/${app_name}/main.go
COPY internal/ internal/

RUN make test
RUN make build

# ----------------------------------- Final -----------------------------------
# Using this instead of scratch with CA certificates and tzdata. This is slightly
# bigger (about 1MB) but makes for a cleaner setup.
FROM gcr.io/distroless/static-debian12:nonroot AS final

ARG artifacts

# Creating and copying a dummy file to avoid buildkit skipping this stage.
# The reason for copying the file here, in stead of the build stage, is that the
# build and lint stage can then be run in parallell.
COPY --from=lint /tmp/dummy.txt /tmp/dummy.txt 

COPY --from=build --chown=nonroot:nonroot /app/${artifacts}/app.bin /app

ENTRYPOINT ["/app"]
