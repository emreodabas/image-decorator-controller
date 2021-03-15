#!/usr/bin/env bash

KUBE_SYSTEM_DEPLOY="test-deployment"
KUBE_SYSTEM_DAEMON="test-daemonset"
TEST_NS="test"


echo "$GREEN >>>> Lets check controller logs for a while $NC"
read -p "$YLW↓ Enter to continue ↓$NC"

timeout 10s kubectl logs -f -c manager -n kube-builder-system $(kubectl get po -n kube-builder-system -o name)

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>Let's check what happened to deployments"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
echo "$RED EXPECTED: kubermatico/nginx$NC"
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>Let's check what happened to daemonsets"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
echo "$RED EXPECTED: kubermatico/nginx$NC"
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>> Deployment under kube-system"
echo ">>kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep \" image:\""
echo "$RED EXPECTED: nginx$NC"
kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep " image:"

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>> Daemonsets under kube-system"
echo ">>kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep \" image:\""
echo "$RED EXPECTED: nginx$NC"
kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep " image:"

echo "$GREEN"
echo ">>>>>>>>>>>>>>>>   DONE   <<<<<<<<<<<<<<<<<<<<"
