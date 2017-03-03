Launch a WebServer (Hello World)

```
kubectl run hello-node --image=kiranmova/hellonode:v1 --port=8080
kubectl get pods
kubectl get deployments
kubectl describe pods hello-node
kubectl get events
```
Wait for the pod to get ready

```
kubectl expose deployment hello-node --type=LoadBalancer
kubectl describe service hello-node
kubectl get nodes
curl k8s-node-1:30220
curl k8s-node-2:30220
```

Upgrade

```
kubectl set image deployment/hello-node hello-node=kiranmova/hellonode:v2
kubectl get pods
curl k8s-node-1:30220
curl k8s-node-2:30220
```


Cleanup
```
kubectl delete service hello-node
kubectl delete deployment hello-node
kubectl get pods
```
