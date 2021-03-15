#!/usr/bin/env bash

KUBE_SYSTEM_DEPLOY="test-deployment"
KUBE_SYSTEM_DAEMON="test-daemonset"
TEST_NS="test"

# if parameter is exist don't delete
if [ -n "$1" ]; then
  ./resetWorkloads.sh "test"
fi

echo ">>>>Showing controller is not deployed"
echo ">>kubectl get deploy -n kube-builder-system"
kubectl get deploy -n kube-builder-system
echo ">>kubectl get deploy -A "
kubectl get deploy -A
echo ">>kubectl get daemonset -A "
kubectl get daemonset -A


read -p "↓↓↓↓↓"

echo ">>>> creating $TEST_NS namespace"
kubectl create ns $TEST_NS

read -p "↓↓↓↓↓"

echo ">>>>Deployment & DaemonSet are creating under $TEST_NS namespace"
echo ">>kubectl create deploy nginx-deploy --image=nginx -n $TEST_NS"
kubectl create deploy nginx-deploy --image=nginx -n $TEST_NS
echo "kubectl apply -f daemonset.yml"
kubectl apply -f daemonset.yml

read -p "↓↓↓↓↓"

echo ">>>>showing images are public "
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep " image:""
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"
echo ">>kubectl get daemonset -o yaml -n $TEST_NS | grep \" image:\""
kubectl get daemonset -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>>Deployment & DaemonSet are creating under kube-system"
echo ">>kubectl create deploy $KUBE_SYSTEM_DEPLOY --image=nginx -n kube-system"
kubectl create deploy $KUBE_SYSTEM_DEPLOY --image=nginx -n kube-system
echo ">>kubectl apply -f daemonset-sys.yml"
kubectl apply -f daemonset-sys.yml

read -p "↓↓↓↓↓"

echo ">>>>Showing images are public"
echo ">>kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep \" image:\""
kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep " image:"
echo ">>kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep \" image:\""
kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>>installing controller to k8s !! need kustomize !!"
## !!Optional
echo ">>make docker-publish -C ../"
make docker-publish -C ../
echo ">>make deploy -C ../"
make deploy -C ../

read -p "↓↓↓↓↓"

echo ">>>>showing controller is deployed"
echo ">>kubectl get deploy -n kube-builder-system"
kubectl get deploy -n kube-builder-system
echo ">>kubectl get po -n kube-builder-system"
kubectl get po -n kube-builder-system

echo ">>>> Lets check controller logs for a while"
kubectl logs -f -c manager -n kube-builder-system $(kubectl get po -n kube-builder-system -o name)

read -p "↓↓↓↓↓"


echo ">>>>Let's check what happened to deployments"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"


echo ">>>>Let's check what happened to daemonsets"
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\""
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>> Deployment under kube-system"
echo ">>kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep \" image:\""
kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep " image:"

read -p "↓↓↓↓↓"

echo ">>>> Daemonsets under kube-system"
echo ">>kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep \" image:\""
kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep " image:"

echo ">>>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<<<<<<<<"
echo ">>>>>>>>>>>>>>>>          <<<<<<<<<<<<<<<<<<<<"
echo ">>>>>>>>>>>>>>>>   DONE   <<<<<<<<<<<<<<<<<<<<"
echo ">>>>>>>>>>>>>>>>          <<<<<<<<<<<<<<<<<<<<"
echo ">>>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<<<<<<<<"
