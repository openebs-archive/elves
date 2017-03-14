#!/bin/bash

#removing older files 

rm -f /tmp/cfssl* && rm -rf /tmp/certs && mkdir -p /tmp/certs

echo "Installing Certs Library CFSSL..."
curl -L https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -o /tmp/cfssl
chmod +x /tmp/cfssl
sudo mv /tmp/cfssl /usr/local/bin/cfssl

curl -L https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -o /tmp/cfssljson
chmod +x /tmp/cfssljson
sudo mv /tmp/cfssljson /usr/local/bin/cfssljson

mkdir -p $HOME/certs
sudo chmod -R a+rw $HOME/certs
cd $HOME/certs

echo "Creating CA certs and key..."
cat > ca-csr.json <<EOF
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "IN",
      "L": "BLR",
      "O": "Kubernetes",
      "OU": "CA",
      "ST": "Karnataka"
    }
  ]
}
EOF

# Generate CA cert and key
cfssl gencert -initca ca-csr.json | cfssljson -bare ca

cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "8760h"
    },
    "profiles": {
      "client": {
        "usages": ["signing", "key encipherment", "server auth", "client auth"],
        "expiry": "8760h"
      }
    }
  }
}
EOF

cat > client-csr.json <<EOF
{
  "CN": "kubernetes",
  "hosts": [
    "k8s-master-1",
    "k8s-master-2",
    "k8s-node-1",
    "k8s-node-2",
    "host-01",
    "host-02",
    "172.28.128.1",
    "172.28.128.2",
    "172.28.128.3",
    "172.28.128.4",
    "172.28.128.5",
    "172.28.128.6",
    "172.28.128.7",
    "172.28.128.8",
    "172.28.128.9",
    "172.28.128.10",
    "172.28.128.11",
    "172.28.128.12",
    "172.28.128.13",
    "172.28.128.14",
    "127.0.0.1"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "IN",
      "L": "BLR",
      "O": "Kubernetes",
      "OU": "Cluster",
      "ST": "Karnataka"
    }
  ]
}
EOF

cfssl gencert \
	  -ca=ca.pem \
	    -ca-key=ca-key.pem \
	      -config=ca-config.json \
	        -profile=client \
		  client-csr.json | cfssljson -bare client
#Renaming certs
mv ca.pem ca.crt
mv ca-key.pem ca.key
mv client.pem client.crt
mv client-key.pem client.key

#Moving certs to config dir. This has to be with flannel script
#sudo mkdir /etc/etcd

#sudo cp $HOME/certs/ca.crt /etc/etcd/
#sudo cp $HOME/certs/client.crt /etc/etcd/
#sudo cp $HOME/certs/client.key /etc/etcd/
