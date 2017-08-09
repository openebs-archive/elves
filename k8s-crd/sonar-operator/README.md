
A sonar device is installed in an submarine (vessel, aka kubernetes cluster) that detects submarines (aka custom-resources) that come in its vicinity. A sonar operator (aka custom-controller) job is to keep looking at the sonar (monitor), get the details of the detected submarines and take appropriate actions. For the sake of this example, let us assume the sonar-operator is trying to raise an alert (aka. publish event), whenever a non-US submarine is detected. 


## Building 
```
make build
```

## Setup CRD
```
kubectl apply -f submarine-crd.yaml
kubectl get crd
```

## Running as Binary
```
sonar-operator -v 4 --kubeconfig=$HOME/.kube/config
```
- -v specify the log level (default is 2 only errors)

## Launching Submarines
```
kubectl apply -f submarine-whale.yaml
kubectl get sub
```

