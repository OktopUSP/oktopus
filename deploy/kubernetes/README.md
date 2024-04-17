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


<b>Node Ports</b>

For this deployment, we are not using a load balancer and kubernetes is deployed on-premises so we are using Nodeports to insource the client traffic into cluster. below the ports set on deployment files:

1. MQTT broker service (mqtt-svc): 30000
2. Frontend (frontend-svc): 30001
3. SocketIO: (socketio-svc): 30002
4. Controller (controller-svc): 30003
5. WebSocket (ws-svc): 30005

Before deploying the files, edit the frontend.yaml file to set the correct enviroment variables:

```yaml
env:
    - name: NEXT_PUBLIC_REST_ENDPOINT
        value: "<FRONTEND_IP>:30003"
    - name: NEXT_PUBLIC_WS_ENDPOINT
        value: "<FRONTEND_IP>:30005"
``` 

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