#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/rebuy-de/kubernetes-pod-restarter/pkg/generated \
github.com/rebuy-de/kubernetes-pod-restarter/pkg/apis \
lifecycle:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
