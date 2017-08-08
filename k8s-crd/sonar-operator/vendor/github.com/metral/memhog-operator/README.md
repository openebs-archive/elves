# Creating Custom Operators

The `memhog-operator` is an example on how to create custom operators for
Kubernetes.

The purpose of the `memhog-operator` is to watch for Pods in a namespace and
monitor its memory usage. If the memory consumption of the Pod crosses a
threshold, it will be vertically autoscaled by the operator.

Specifically, the operator will deploy a new copy of the Pod with a higher
set of resource requests and limit, and then terminate the original Pod.
The details of the higher resources are held within an `AppMonitor`,
a [CustomResourceDefinition
(CRD)](https://github.com/kubernetes/apiextensions-apiserver/blob/fbe70034cb9becd97bf8b6207f918c73cadd330e/pkg/apis/apiextensions/types.go#L119-L129).

[memhog](https://github.com/metral/memhog): An example Pod that this operator would monitor.

> Note: The `memhog-operator` is strictly for educational & demo purposes. It is not intended
to be used for any other use-cases.

## Operator Structure

The `memhog-operator` is a combination of a CustomResourceDefinition
known as the `AppMonitor`, and a custom Controller to enforce state.

The `AppMonitor` encapsulates the autoscaling details for a Pod.
The controller watches a Namespace for an `AppMonitor`, and for Pods that wish 
to be monitored (via Annotation). It then applies the operational
thresholds and requirements declared in the `AppMonitor` onto the Pod.

## Requirements
* Kubernetes v1.7.0+
* Prometheus on k8s with cAdvisor exposing `container_memory_usage_bytes`
* glide v0.11.1

## Process

* To monitor the Pod's resource memory consumption, the operator requires that the Pod have an annotation in its `spec.template.metadata` to associate itself with the `memhog-operator`.

  e.g. The [memhog](https://github.com/metral/memhog) example is annotated as such:

  <pre><code>apiVersion: extensions/v1beta1
  kind: Deployment
  metadata:
    name: memhog
  spec:
    replicas: 1
    template:
      metadata:
        labels:
          name: memhog
        <b>annotations:
          app-monitor.kubedemo.com/monitor: "true"</b>
      spec:
        containers:
        - name: memhog
          image: quay.io/metral/memhog:v0.0.1
          imagePullPolicy: Always
          <b>resources:
            limits:
              memory: 384Mi
            requests:
              memory: 256Mi</b>
          ...
  </code></pre>

* Run the annotated Pod, noting the resource requests & limits set that are
also required by the Operator.

* In the Pod's namespace there must also be an instantiated object of the custom
`AppMonitor` CRD that the Operator depends on. For example, this `AppMonitor`
states that any Pod being monitored by the `memhog-operator` will have its
resources doubled when 75% or more of its memory has been used e.g.:

  ```
  apiVersion: kubedemo.com/v1
  kind: AppMonitor
  metadata:
    name: johnny-cache
  spec:
    memThresholdPercent: 75   # Percentage of (memory used) / (memory limit)
    memMultiplier: 2          # Multiplier factor used to increase memory resource requests & limits
  ```
* The `memhog-operator` will watch the Pod's memory usage by querying
Prometheus, applying the `AppMonitor` to the Pod as memory usage is retrieved.  If the `AppMonitor` threshold is crossed,
the Operator will redeploy the Pod with higher resource requests & limits based
on the multiplier.
* If the Pod is redeployed, it will have updated and increased resource requests & limits
e.g.:
```
  ...
  resources:
    limits:
      memory: 768Mi
    requests:
      memory: 512Mi
  ...
```

### Building & Running

```
// Build
$ glide install -s -u -v
$ make
```


#### Run Locally
> Note: Prometheus is assumed to be running locally on http://localhost:9090

```
$ $GOPATH/bin/memhog-operator -v2 --prometheus-addr=http://localhost:9090 --kubeconfig=$HOME/.kube/config
```

#### Run on Kubernetes
> Note: Prometheus is assumed to be running by default on http://prometheus.tectonic-system:9090 (e.g. Tectonic Cluster) by the Operator Deployment

```
// Create cluster role & cluster role binding to work with CRD's.
$ kubectl create -f k8s/roles/role.yaml

$ kubectl create -f k8s/deploy/memhog-operator-deploy.yaml
```

### References

* [ThirdPartyResources (TPR) -> CustomResourceDefinitions (CRD)](https://groups.google.com/forum/#!msg/kubernetes-dev/R749_-L_ssc/p-3tyl6mAQAJ)
* [CRD Example](https://github.com/kubernetes/apiextensions-apiserver/tree/master/examples/client-go)
* [The TPR is dead. Long live the CRD: Custom Resources in Kubernetes v1.7](https://coreos.com/blog/custom-resource-kubernetes-v17)
* [Working with Controllers](https://github.com/kubernetes/community/blob/master/contributors/devel/controllers.md)
* [Writing a Custom Controller - Aaron Levy (CoreOS)](https://github.com/aaronlevy/kube-controller-demo
)
  * [KubeCon EU 2017 Video](youtu.be/_BuqPMlXfpE?list=PLj6h78yzYM2PAavlbv0iZkod4IVh_iGqV)
* [DaemonSet Controller Code](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/daemon/daemoncontroller.go)

* Deprecated
  * [TPRs are Deprecated](https://github.com/kubernetes/client-go/tree/df46f7f13b3da19b90b8b4f0d18b8adc6fbf28dc/examples/third-party-resources-deprecated
)
  * [Extend the Kubernetes API with TPRs](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-third-party-resource/)
