#!/usr/bin/env bash

KUBE_SYSTEM_DEPLOY="test-deployment"
KUBE_SYSTEM_DAEMON="test-daemonset"
TEST_NS="test"


echo ">>>>Let's check what happened to deployments"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>>Let's check what happened to daemonsets"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>> Deployment under kube-system"
echo ">>kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n $TEST_NS | grep \" image:\""
kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>> Daemonsets under kube-system"
echo ">>kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep \" image:\""
kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep " image:"
