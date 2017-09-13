kubectl create -f service-server.yaml
kubectl create -f deployment-server.yaml

kubectl create -f prometheus.yaml
kubectl create -f deployment-prometheus.yaml
kubectl create -f service-prometheus.yaml

kubectl get pods
kubectl get svc
kubectl get deployment
kubectl describe confimap
