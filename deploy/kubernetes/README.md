# Oktopus Kubernetes

## Requirements

Kubernetes 1.28+

### Standalone Installation

Single Node:
* 8 vCPUs
* 8 GB RAM

# Installation

## Download Files

```shell
git clone https://github.com/OktopUSP/oktopus
export DEPLOYMENT_PATH=oktopus/deploy/kubernetes
```

## HAProxy Ingress Controller

```shell
helm install haproxy-kubernetes-ingress haproxytech/kubernetes-ingress \
  --create-namespace \
  --namespace haproxy-controller \
  --set controller.kind=DaemonSet \
  --set controller.daemonset.useHostPort=true
```

## MongoBD

```shell
# Mongo DB Operator at mongodb namespace
helm repo add mongodb https://mongodb.github.io/helm-charts

helm install community-operator mongodb/community-operator --namespace mongodb --create-namespace

# Mongo DB ReplicaSet
export DEPLOYMENT_PATH=oktopus/deploy/kubernetes

kubectl apply -f $DEPLOYMENT_PATH/mongodb.yaml -n mongodb

# Check Installation
kubectl get pods -n mongodb
```

## NATS Server

```shell
# Download the NATS charts
helm repo add nats https://nats-io.github.io/k8s/helm/charts/

# Install NATS with Jetstream Enabled
helm install nats nats/nats --set config.jetstream.enabled=true
```
 
## Oktopus

```shell
kubectl apply -f $DEPLOYMENT_PATH/mqtt.yaml
kubectl apply -f $DEPLOYMENT_PATH/mqtt-adapter.yaml
kubectl apply -f $DEPLOYMENT_PATH/adapter.yaml
kubectl apply -f $DEPLOYMENT_PATH/controller.yaml
kubectl apply -f $DEPLOYMENT_PATH/socketio.yaml
kubectl apply -f $DEPLOYMENT_PATH/frontend.yaml
kubectl apply -f $DEPLOYMENT_PATH/ws.yaml
kubectl apply -f $DEPLOYMENT_PATH/ws-adapter.yaml
```

### Checking cluster status:

```shell

kubectl get pods
kubectl get svc

```
