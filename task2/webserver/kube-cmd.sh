#!/bin/bash
kubectl create -f volume.yaml
kubectl create -f volume-claim.yaml
kubectl run webserver --image="webserver:v1" --overrides="$(cat override.json)"
kubectl autoscale deployment webserver --cpu-percent=50 --min=1 --max=10
kubectl expose deployment/webserver --type=LoadBalancer
