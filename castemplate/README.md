Upgrade of cstor pool has following steps -<br>
```
1. Patch CSP cr.
2. Patch SP cr.
3. Patch pool deployment.
4. Wait for patch status.
5. Delete older replicaset of Pool deploymnet.
6. Perform post upgrade activity.(exec into pool pod and run some command).
```
These are the runtask to upgrade cstor pool- <br>
```
1. task-getcspdeploy (this task provides pool name, namespace, claimname)
2. task-getcspcr (this task provides csp uuid)
3. task-patchcsp (this runtask patches csp cr)
4. task-patchsp (this runtask patches sp cr)
5. task-listrs (this runtask lists older replicaset of that pool deploymnet)
6. task-patchcspdeploy (this runtask patches pool deployment)
7. task-deleters (this runtask deletes older replicaset for that pool deployment)
8. task-postupgradejob (this runtask performs some post upgrde task)
```
Upgrade stape -

1. list out all spc `kubectl get spc`
2. form that spc list out pool deployments,csp,sp
```
kubectl get deploy -n opnenes -l app=cstor-pool,openebs.io/storage-pool-claim=<claim name>
kubectl get sp , csp -l openebs.io/storage-pool-claim=<claim name>
```
3. select one pool and update it's name in `task-getcspdeploy` runtask.
In step 1 castemplate task's will be-
```
    tasks:
    - task-getcspdeploy
    - task-getcspcr
    - task-patchcsp
    - task-patchsp
    - task-listrs
    - task-patchcspdeploy
    - task-deleters
```
In step 2 (when new pods are in running state) castemplate task's will be-
```
    tasks:
    - task-getcspdeploy
    - task-postupgradejob
```
