#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
KUBECONFIG_PATH="${KUBECONFIG:-$SCRIPT_DIR/../.omc/plans/kubeconfig}"
NAMESPACE="${NAMESPACE:-zcicd}"
RELEASE_NAME="${RELEASE_NAME:-zcicd}"
VALUES_FILE="${VALUES_FILE:-$SCRIPT_DIR/helm/values-dev.yaml}"

export KUBECONFIG="$KUBECONFIG_PATH"

echo "==> Using kubeconfig: $KUBECONFIG_PATH"
echo "==> Namespace: $NAMESPACE"

# Create namespace
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# Install/upgrade with Helm
helm upgrade --install "$RELEASE_NAME" \
  "$SCRIPT_DIR/helm/zcicd" \
  -n "$NAMESPACE" \
  -f "$VALUES_FILE" \
  --wait --timeout 10m

echo "==> Deployment complete. Checking pods:"
kubectl get pods -n "$NAMESPACE"
