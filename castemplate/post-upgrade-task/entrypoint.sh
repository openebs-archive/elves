#!/bin/bash

set -ex

if [ -z "$ns" ]; then
    echo "ERROR: namespace not present"
    exit 1
fi

if [ -z "$spc" ]; then
    echo "ERROR: namespace not present"
    exit 1
fi

if [ -z "$csp" ]; then
    echo "ERROR: namespace not present"
    exit 1
fi

pool_rs=$(kubectl get replicaset -n $ns \
    -l app=cstor-pool,openebs.io/storage-pool-claim=$spc \
    -o jsonpath="{range .items[?(@.metadata.ownerReferences[0].name== '$csp')]}{@.metadata.name}{end}")

if [ -z "$pool_rs" ]; then
    echo "unable to find pool replicaset"
    exit 1
fi

pool_pod=$(kubectl get pod -n $ns\
    -l app=cstor-pool,openebs.io/storage-pool-claim=$spc \
    -o jsonpath="{range .items[?(@.metadata.ownerReferences[0].name== '$pool_rs')]}{@.metadata.name} {end}")

if [ -z "$pool_pod" ]; then
    echo "unable to find pool pod"
    exit 1
fi

pool_name=""
cstor_uid=""
cstor_uid=$(kubectl get pod $pool_pod -n $ns \
    -o jsonpath="{.spec.containers[*].env[?(@.name=='OPENEBS_IO_CSTOR_ID')].value}" | awk '{print $1}')
pool_name="cstor-$cstor_uid"
quorum_set=$(kubectl exec $pool_pod -n $ns -c cstor-pool-mgmt -- zfs set quorum=on $pool_name)
rc=$?
if [[ ($rc -ne 0) ]]; then
    echo "Error: failed to set quorum for pool $pool_name"
    exit 1
fi
output=$(kubectl exec $pool_pod -n $ns -c cstor-pool-mgmt -- zfs get quorum)
rc=$?
if [ $rc -ne 0 ]; then
    echo "ERROR: while executing zfs get quorum for pool $pool_name, error: $rc"
    exit 1
fi
no_of_non_quorum_vol=$(echo $output | grep -wo off | wc -l)
if [ $no_of_non_quorum_vol -ne 0 ]; then
    echo "Few($no_of_non_quorum_vol) of quorum values are having inappropriate values for quorum"
    exit 1
fi

