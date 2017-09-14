## prometheus-monitoring

This contains very basic server program and it's endpoint monitoring with prometheus. Currently this setup only works on minikube.This repo is highly motivated from [this repo](https://github.com/marselester/prometheus-on-kubernetes). rbac-setup file added to run this setup on k8s cluster.

# Steps for minikube
1. cd to this folder
2. run `minikube start` ([minikube](https://github.com/kubernetes/minikube) is required)
3. run script `sh server-script.sh`
4. check pods `kubectl get pods`
5. check svc `kubectl get svc`
6. check configmap `kubectl get configmap`
7. Once all working fine, run `minikube service goserver-service` (It will start server on a browser)
8. Now run `minikube service prometheus-service` (It will launch prometheus expression browser, where you can monitor your app)

# Steps for kubeadm
1. cd to this folder
2. run `kubectl create -f rbac-setup.yaml`
3. repeat step 3 to 6 like in minikube
4. once you get status of containers runnig run `ssh -NL 1234:localhost:32514 ubuntu@<ip_address of Node>` and open localhost:1234 in your browser to access ui of prometheus from outside k8s cluster setup.
5. run `ssh -NL 1235:localhost:32500 ubuntu@<ip_address of Node>` and open localhost:1235 in your browser to hit at the /latest/volumes from outside k8s cluster setup.
6. Now you will be able to see the targets in status (dropdown button) and custom metrics.
