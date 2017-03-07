
```
sudo mkdir -p /var/lib/kubernetes
sudo cp ../certs/ca.crt /var/lib/kubernetes
sudo cp ../certs/client.crt /var/lib/kubernetes
sudo cp ../certs/client.key /var/lib/kubernetes


sudo cp authorization-policy.jsonl /var/lib/kubernetes/
sudo cp token.csv /var/lib/kubernetes/

Update the IP address of the etcd service in the kube-apiserver.service. 

sudo cp kube-apiserver.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start kube-apiserver
sudo systemctl status kube-apiserver
tail -f /var/log/syslog 

sudo cp kube-controller-manager.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start kube-controller-manager
sudo systemctl status kube-controller-manager
tail -f /var/log/syslog 

sudo cp kube-scheduler.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl start kube-scheduler
sudo systemctl status kube-scheduler
tail -f /var/log/syslog 

kubectl get componentstatus
kubectl get componentstatuses

sudo systemctl enable kube-apiserver
sudo systemctl enable kube-controller-manager
sudo systemctl enable kube-scheduler
```

