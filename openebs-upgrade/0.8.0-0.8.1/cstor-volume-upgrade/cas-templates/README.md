# Cstor Volume Upgrade Using CAS-Template (0.8.0 -> 0.8.1)

After Cstor pool upgrade, now we will upgrade the volume from version 0.8.0 to 0.8.1.

Steps we will be going through for upgrading a volume are  -

(Note: We are using CAS-Templates for various upgrade steps (every upgrade step will be a runtask))

1. Get CASType, storage class name and storage class namespace for volume given (For now, this
   cas-template is for upgrading cstor volumes only).

cas_type=`kubectl get pv $pv -o jsonpath="{.metadata.annotations.openebs\.io/cas-type}"`
=> cstor

sc_ns=`kubectl get pv pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4 -o jsonpath="{.spec.claimRef.namespace}"`
=> default

sc_name=`kubectl get pv pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4 -o jsonpath="{.spec.storageClassName}"`
=> openebs-cstor-disk

- RunTask being used (get-vol-info)

2. Get resource version for the storage class being used by volume.

sc_res_ver=`kubectl get sc openebs-cstor-disk -n default -o jsonpath="{.metadata.resourceVersion}"`
=> 4102481

- RunTask being used (get-sc-res-version)

3. Get target-deployment, target-svc, cstor_volumes, cstor_volume_replicas(cvr), target_old_rs..

### STEP: Generate deploy, replicaset and container names from PV
#### NOTES: Ex: If PV="pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4"
####  then c-dep: pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-target

c_dep=$(kubectl get deploy -n openebs -l openebs.io/persistent-volume=pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4,openebs.io/target=cstor-target -o jsonpath="{.items[*].metadata.name}")
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-target

- RunTask being used (list-ctrl-deployment)

c_svc=$(kubectl get svc -n openebs -l openebs.io/persistent-volume=pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4,openebs.io/target-service=cstor-target-svc -o jsonpath="{.items[*].metadata.name}")
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4

- RunTask being used (list-ctrl-svc)

c_vol=$(kubectl get cstorvolumes -l openebs.io/persistent-volume=pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4 -n openebs -o jsonpath="{.items[*].metadata.name}")
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4

- RunTask being used (list-cstor-volumes-cr)

c_replicas=$(kubectl get cvr -n openebs -l openebs.io/persistent-volume=pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4 -o jsonpath="{range .items[*]}{@.metadata.name};{end}" | tr ";" "\n")
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-demo-cstor-pool-cc9d
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-demo-cstor-pool-ssa7
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-demo-cstor-pool-zyq5

- RunTask being used (list-cstor-volume-replicas)

### Fetch the older target and replica - ReplicaSet objects which need to be deleted before upgrading. 
### If not deleted, the new pods will be stuck in creating state - due to affinity rules.

c_rs=$(kubectl get rs -n openebs -o name -l openebs.io/persistent-volume=pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4 | cut -d '/' -f 2)
=> pvc-04a7fec6-3a6b-11e9-a409-42010a8002d4-target-799f577f68

- RunTask being used (list-target-old-rs)

4. Patch volume target (deployment) strategy to "Recreate".
    - RunTask being used (patch-deployment-recreate-strategy).

5. Patch volume target deployment with the changes for 0.8.1 version i.e. container images, annotations,
   labels, etc.
    - RunTask being used (patch-target-deployment)

6. Delete old replicaset for target deployment.
    - RunTask being used (delete-target-old-replicaset)

7. Patch target svc
    - RunTask being used (patch-target-svc)

8. Patch Cstor volume CR
    - RunTask being used (patch-cstor-volume-cr)

9.  Patch Cstor Volume replicas
    - RunTask being used (patch-cstor-volume-replica)
