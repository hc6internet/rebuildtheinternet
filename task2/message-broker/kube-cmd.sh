#!/bin/bash
kubectl run message-broker --image="rabbitmq:v1" --overrides="$(cat override.json)"
kubectl expose deployment/message-broker
