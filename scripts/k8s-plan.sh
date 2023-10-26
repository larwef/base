#!/usr/bin/env bash

set -euxo pipefail

kustomize build deployments/kubernetes | envsubst >"$MANIFEST_OUTPUT"
# To avoid failing when there is a diff.
set +e
kubectl --context "$K8S_CONTEXT" ${K8S_NAMESPACE:+-n="$K8S_NAMESPACE"} diff -f "$MANIFEST_OUTPUT"
exit 0 # Because make will take it as an error if there is a diff.
