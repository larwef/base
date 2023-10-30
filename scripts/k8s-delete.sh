#!/usr/bin/env bash

set -euxo pipefail

kubectl --context "$K8S_CONTEXT" ${K8S_NAMESPACE:+-n="$K8S_NAMESPACE"} delete -f "$MANIFEST_INPUT"
