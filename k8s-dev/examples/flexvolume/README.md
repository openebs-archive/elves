On each of the minions:

```
sudo mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/cb~temp
sudo cp temp /usr/libexec/kubernetes/kubelet-plugins/volume/exec/cb~temp
sudo chmod +x /usr/libexec/kubernetes/kubelet-plugins/volume/exec/cb~temp
```

Restart kubelet

```
sudo systemctl restart kubelet
sudo systemctl status kubelet
```

Check the cb/temp is loaded

```
grep "Loaded volume plugin" /var/log/syslog
```
