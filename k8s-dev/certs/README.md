
Install cfssl 
```
sh install_cfssl.sh
```

Use the default ca-csr.json provided in the ./certs/ca-csr.json 

Generate the CA certs
```
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
```

Verify the CA Cert
```
openssl x509 -in ca.pem -text -noout
```

Create Client Certificate - A generic one for all the nodes using the ./certs/client-csr.json

```
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client-csr.json |  cfssljson -bare client
```

Rename

```
mv ca.pem ca.crt
mv ca-key.pem ca.key
mv client.pem client.crt
mv client-key.pem client.key
```

