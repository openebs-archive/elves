## prometheus-monitoring

This contains very basic server program and it's endpoint monitoring with prometheus. Currently this setup only works on minikube.

# Steps
* cd to this folder
* run `minikube start` ([minikube](https://github.com/kubernetes/minikube) is required)
* run script `sh server-script.sh`
* check pods `kubectl get pods`
* check svc `kubectl get svc`
* check configmap `kubectl get configmap`
* Once all working fine, run `minikube service goserver-service` (It will start server on a browser)
* Now run `minikube service prometheus-service` (It will launch prometheus expression browser, where you can monitor your app)
