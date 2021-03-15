#!/usr/bin/env bash

KUBE_SYSTEM_DEPLOY="test-deployment"
KUBE_SYSTEM_DAEMON="test-daemonset"
TEST_NS="test"
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
NC=$(tput sgr0)
YLW=$(tput setaf 3)
BLUE=$(tput setaf 4)

# if parameter is exist don't delete
if [ -n "$1" ]; then
  ./resetWorkloads.sh "test"
fi

echo "$GREEN >>>>Showing controller is not deployed"
echo ">>kubectl get deploy -n kube-builder-system$NC"
kubectl get deploy -n kube-builder-system
echo "$GREEN >>kubectl get deploy -A $NC"
kubectl get deploy -A
echo "$GREEN >>kubectl get daemonset -A $NC"
kubectl get daemonset -A

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>> creating $TEST_NS namespace$NC"
kubectl create ns $TEST_NS

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>Deployment & DaemonSet are creating under $TEST_NS namespace"
echo ">>kubectl create deploy nginx-deploy --image=nginx -n $TEST_NS $NC"
kubectl create deploy nginx-deploy --image=nginx -n $TEST_NS
echo "$GREEN kubectl apply -f daemonset.yml $NC"
kubectl apply -f daemonset.yml

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>showing images are public "
echo ">>kubectl get deploy -o yaml -n $TEST_NS | grep \" image:\" $NC"
kubectl get deploy -o yaml -n $TEST_NS | grep " image:"
echo "$GREEN >>kubectl get daemonset -o yaml -n $TEST_NS | grep \" image:\"$NC"
kubectl get daemonset -o yaml -n $TEST_NS | grep " image:"

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>Deployment & DaemonSet are creating under kube-system"
echo ">>kubectl create deploy $KUBE_SYSTEM_DEPLOY --image=nginx -n kube-system$NC"
kubectl create deploy $KUBE_SYSTEM_DEPLOY --image=nginx -n kube-system
echo "$GREEN >>kubectl apply -f daemonset-sys.yml $NC"
kubectl apply -f daemonset-sys.yml

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>Showing images are public"
echo ">>kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep \" image:\"$NC"
kubectl get deploy $KUBE_SYSTEM_DEPLOY -o yaml -n kube-system | grep " image:"
echo "$GREEN >>kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep \" image:\"$NC"
kubectl get daemonset $KUBE_SYSTEM_DAEMON -o yaml -n kube-system | grep " image:"

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>installing controller to k8s !! need kustomize !!$NC"
## !!Optional
echo "$GREEN >>make docker-publish -C ../$NC"
make docker-publish -C ../
echo "$GREEN >>make deploy -C ../ $NC"
make deploy -C ../

read -p "$YLW↓ Enter to continue ↓$NC"

echo "$GREEN >>>>showing controller is deployed"
echo ">>kubectl get deploy -n kube-builder-system$NC"
kubectl get deploy -n kube-builder-system
echo "$GREEN >>kubectl get po -n kube-builder-system$NC"
kubectl get po -n kube-builder-system

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
