#!/bin/bash

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/cloud/deploy.yaml > /dev/null
kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission

if ! docker image inspect custom/element:v1.0.0 > /dev/null
then
  git clone https://github.com/vector-im/element-web.git
  docker build -t custom/element:v1.0.0 ./element-web/
  rm -rf element-web
fi

# clean up
helm uninstall element > /dev/null
helm uninstall synapse > /dev/null
helm uninstall dendrite > /dev/null

# deploy services
helm install synapse "./synapse/helm/"
helm install element "./element/helm/"
