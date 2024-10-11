# sqlite

Delete this folder if you don't need it. Or remember to hook it up if you're
using it.

Note that you have to use cgo to build with this sqlite setup. If you're going
to cross compile it might get a bit difficult. This is an example which should
work when setting

```bash
GOOS=linux \
GOARCH=arm \
CGO_ENABLED=1 \
CC=arm-linux-gnueabihf-gcc \
CXX=arm-linux-gnueabihf-g++ \
```

```dockerfile

FROM golang:1.22 as build

ARG app_name

# Needed for cross compiling for arm.
RUN apt-get update && apt-get install -y \
  gcc-arm-linux-gnueabihf \
  g++-arm-linux-gnueabihf

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

```

