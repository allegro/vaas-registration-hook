#!/usr/bin/env bash

set -e

sudo mkdir /hooks
sudo cp target/vaas-hook /hooks

echo "Starting myservice-pod..."
kubectl create -f examples/service-with-lifecycle.yaml

until kubectl get pods | tee | grep Running
do
  echo "Waiting for myservice-pod to start successfully..."
  kubectl describe pod myservice-pod
  sleep 30
done

RESULT=$(kubectl exec -it myservice-pod -- curl -v --fail http://localhost:80/v1/catalog/service/myservice)
echo $RESULT | grep myservice
echo $RESULT | grep "k8sPodNamespace: default"
