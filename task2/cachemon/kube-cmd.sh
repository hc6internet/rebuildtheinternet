#!/bin/bash
kubectl run cachemon --image=cachemon:v1 --port=9000 --image-pull-policy=Never \
    --env="MON_PORT=9000" \
    --env="RMQ_HOST=message-broker"
kubectl expose deployment/cachemon --type=LoadBalancer
