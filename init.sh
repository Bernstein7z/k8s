#!/bin/bash

# start helm
helm uninstall element synapse > /dev/null
helm install synapse "./synapse/helm/"
helm install element "./element/helm/"

# start k8s dashboard
pkill -9 -f "kubectl proxy"
kubectl proxy &
