```
sudo mkdir /etc/etcd
sudo cp ../certs/ca.crt /etc/etcd/
sudo cp ../certs/client.crt /etc/etcd/
sudo cp ../certs/client.key /etc/etcd/

sudo chown -R etcd /etc/etcd/
sudo chgrp -R etcd /etc/etcd/
```

Set etcd configuration (/etc/default/etcd) and change the gray values appropriately:

Verify the hostname of master ( k8s-master-1) and the IP address used in all the properties

```
ETCD_NAME=k8s-master-1
ETCD_DATA_DIR="/var/lib/etcd"
ETCD_LISTEN_PEER_URLS="https://172.28.128.10:2380"
ETCD_LISTEN_CLIENT_URLS="https://172.28.128.10:2379,http://localhost:2379"
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://172.28.128.10:2380"
ETCD_INITIAL_CLUSTER="k8s-master-1=https://172.28.128.10:2380"
ETCD_INITIAL_CLUSTER_STATE="new"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster-0"
ETCD_ADVERTISE_CLIENT_URLS="https://172.28.128.10:2379"
ETCD_CERT_FILE="/etc/etcd/client.crt"
ETCD_KEY_FILE="/etc/etcd/client.key"
ETCD_TRUSTED_CA_FILE="/etc/etcd/ca.crt"
ETCD_PEER_CERT_FILE="/etc/etcd/client.crt"
ETCD_PEER_KEY_FILE="/etc/etcd/client.key"
ETCD_PEER_TRUSTED_CA_FILE="/etc/etcd/ca.crt"
```

Restart and verify etcd
```
$ systemctl restart etcd
$ systemctl status etcd
$ etcdctl --ca-file=/etc/etcd/ca.crt cluster-health
```
