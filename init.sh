#!/bin/bash

# start helm
helm uninstall element dendrite > /dev/null
helm install dendrite "./dendrite/helm/"
helm install element "./element/helm/"

# start k8s dashboard
pkill -9 -f "kubectl proxy"
kubectl proxy
