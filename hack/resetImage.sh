#!/usr/bin/env bash

echo "set image name to global one"
kubectl set image deployment/nginx  nginx=nginx

echo "show image name is updated"
kubectl get deploy -o yaml | grep " image:"