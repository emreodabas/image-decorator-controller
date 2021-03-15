#!/usr/bin/env bash

KUBE_SYSTEM_DEPLOY="test-deployment"
KUBE_SYSTEM_DAEMON="test-daemonset"
NS="test"

echo ">> Controller NS is deleting"
kubectl delete ns kube-builder-system  --grace-period=0

echo ">> $NS namespace is deleting"
kubectl delete ns $NS --grace-period=0

echo ">> Test Deployment & DaemonSet are deleting under KUBE-SYSTEM "
kubectl delete deployment $KUBE_SYSTEM_DEPLOY -n kube-system
kubectl delete -f daemonset-sys.yml
