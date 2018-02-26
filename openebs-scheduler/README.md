## OpenEBS Scheduler Design
OpenEBS scheduler is the controller that tries to assign appropriate Kubernetes node to a openebs replica pod. It assigns policies to let the replica pod remain bound to the node where the replica pod got placed.


### OpenEBS Scheduler via Kubernetes Initializer

#### Notes
- Initializers must have a unique fully qualified name. e.g.:
  - initializer.openebs.io
  - replica.initializer.openebs.io
  - jivareplica.initializer.openebs.io
- Initializers should be deployed using a Deployment for:
  - easy upgrades, and
  - auto restarts
- Initializers should explicitly set the list of pending initializers to:
  - exclude itself, or
  - an empty array, 
  - NOTE: This avoids the initializer from getting stuck waiting to initialize
  - Examples:
```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  initializers:
    pending: []
```
- Limit the scope of objects to be initialized to the smallest subset possible using an InitializerConfiguration
  - Examples:
```yaml
apiVersion: admissionregistration.k8s.io/v1alpha1
kind: InitializerConfiguration
metadata:
  name: openebs-scheduler
initializers:
  - name: jivareplica.initializer.openebs.io
    rules:
      - apiGroups:
          - "*"
        apiVersions:
          - "*"
        resources:
          - deployments
```
- Use annotations to enable opting in or out of initialization. Examples:
```
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  annotations:
    "jivareplica.initializer.openebs.io": "true"
    "jivareplica.initializer.openebs.io/version": "v1alpha1"
  name: jiva-volume-replica
...
```

### OpenEBS Scheduler

The openebs scheduler is a [Kubernetes Initializer](https://kubernetes.io/docs/admin/extensible-admission-controllers/#what-are-initializers) that sets [nodeSelector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) to openebs replica deployment.

#### Pseudo Scheduler Logic for OpenEBS Replica 
- for each replica pod's <node-name>
 - kubectl label nodes <node-name> openebs.io/scheduler=<replica-deployment-name>
- end of loop
- get replica deployment
- make a copy of replica deployment
- update the copy to include nodeSelector settings
- kubectl patch replica deployment

### References
- https://github.com/kelseyhightower/kubernetes-initializer-tutorial/
