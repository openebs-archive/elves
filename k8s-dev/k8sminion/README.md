
sudo mkdir -p /var/lib/kubernetes
sudo cp ../certs/ca.crt /var/lib/kubernetes
sudo cp ../certs/client.crt /var/lib/kubernetes
sudo cp ../certs/client.key /var/lib/kubernetes


sudo mkdir /opt/cni
wget https://storage.googleapis.com/kubernetes-release/network-plugins/cni-07a8a28637e97b22eb8dfe710eeae1344f69d16e.tar.gz
sudo tar -xvf cni-07a8a28637e97b22eb8dfe710eeae1344f69d16e.tar.gz -C /opt/cni

sudo mkdir -p /var/lib/kubernetes/cni/net.d/
sudo cp 10-flannel.conf /var/lib/kubernetes/cni/net.d/
sudo cp docker_opts_cni.env /var/lib/kubernetes/cni/docker_opts_cni.env


wget https://github.com/coreos/flannel/releases/download/v0.6.2/flanneld-amd64
chmod +x flanneld-amd64
sudo mv flanneld-amd64 /usr/local/bin/flanneld


(Run this on K8s-master)
etcdctl --ca-file=/var/lib/kubernetes/ca.crt set /coreos.com/network/config '{"Network": "10.200.0.0/16", "SubnetLen":24, "Backend": {"Type": "vxlan"}}'
Add the k8s-master-1 IP address in the /etc/hosts of the minion


sudo cp flannel.service /etc/systemd/system/flannel.service
sudo systemctl daemon-reload
sudo systemctl restart flannel
sudo systemctl status flannel
ip addr show


sudo wget https://get.docker.com/builds/Linux/x86_64/docker-latest.tgz
sudo tar -xvf docker-latest.tgz
sudo cp docker/docker* /usr/bin/

sudo cp docker.service /etc/systemd/system/docker.service
sudo systemctl daemon-reload
sudo systemctl start docker
sudo systemctl status docker
sudo docker ps


wget https://storage.googleapis.com/kubernetes-release/release/v1.5.1/bin/linux/amd64/kubelet
wget https://storage.googleapis.com/kubernetes-release/release/v1.5.1/bin/linux/amd64/kube-proxy
chmod +x kubelet kube-proxy
sudo mv kubelet kube-proxy /usr/local/bin/
sudo mkdir -p /var/lib/kubelet/ 
sudo cp kubeconfig /var/lib/kubelet/

sudo cp kubelet.service /etc/systemd/system/kubelet.service
sudo systemctl daemon-reload
sudo systemctl start kubelet
sudo systemctl status kubelet

sudo cp kubelet.service /etc/systemd/system/kube-proxy.service
sudo systemctl daemon-reload
sudo systemctl start kube-proxy
sudo systemctl status kube-proxy
