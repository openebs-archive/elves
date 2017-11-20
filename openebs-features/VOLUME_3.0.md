## Volume 3.0

> This is the third iteration w.r.t volume provisioning implementation

```yaml
Version: 3.0.0
Owners: @amitkumardas
Github Repo: 
Start Date: 20 Nov 2017
End Date: 
```

#### Motivation

[A]
We have Maya api server caterings to the needs of openebs volume provisioning.
While it has sufficed most of the needs for provisioning, there is good scope
for improvement when it comes to provisioning storage desired by an operator.
There is no one size fits all for an application. The storage characteristics
changes depending on the application that wants to consume the former. The 
application along with its workload patterns & deployment environment 
necessitates tunings to its underlying storage. This is required for the application
to behave similarly even with the changes to load & deployment environment.

We at OpenEBS belive, these kind of requirements can be handled when thought in 
terms of storage policies. A storage policy may be composed of various fine 
granular storage policies.

[B]
Kubernetes has been coming up with new patterns for building controllers. Controllers
are logic that is tied to Kubernetes APIs one one hand & application's operational
logic on the other. It seems good to implement volume provisioning logic as a K8s
controller & hence reap the benefits of K8s lifecycle, non-functional, etc.
features by default. This also assumes, maya to deviate from its earlier philosophy
of being an abstraction to various container orchestrators. In other words, maya
will be tightly coupled with Kubernetes.

[C]
Maya api server needs to cater to various volume types i.e. jiva, cstor, etc
with ease. Any new volume type should be implemented with ease. Maya api server
should establish a pattern that can be easily followed to add or update any volume
types.

[D]
Declarative Code is the new norm. While declarative code has its own rough edges,
we are seeing less of those cases in Kubernetes world. The reason behind this is
the way declarative yamls are mapped with Go's structs. Kubernetes avoids the 
templating stuff in these declarations and instead transforms to strict
types for all its Kind. It has not failed in its attempt so far considering the
popularity of these yamls among the developers as well as operators and testers.
Hence, it makes sense for maya to make use of these yamls as much as possible.

#### Design -- Volume API/Structure

```go
type Volume struct {
  // Name of this volume
  Name        string
  
  // Capacity of this volume
  Capacity    string
  
  // Properties of this volume
  Properties  Properties
  
  // Policies of this volume
  Policies    Policies
}

// Properties holds various references using which volume's properties
// can be determined
type Properties struct {
  // ConfigMap is the name of the ConfigMap that holds volume's properties
  ConfigMap   string
  
  // EndPoint is the url of the service that provides volume's properties
  EndPoint    string
  
  // Priorities contains the prioritized list of references
  // e.g. {"ConfigMap", "EndPoint"} implies attempt fetching of properties from 
  // ConfigMap attribute & if it fails attempt fetching of properties from
  // EndPoint attribute
  Priorities  []VolumeReference
}

// Policies holds various references using which volume's policies
// can be determined
type Policies struct {
  // ConfigMap is the name of the ConfigMap that holds volume's policies
  ConfigMap   string
  
  // EndPoint is the url of the service that provides volume's policies
  EndPoint    string
  
  // Priorities contains the prioritized list of references
  // e.g. {"ConfigMap", "EndPoint"} implies attempt fetching of policies from 
  // ConfigMap attribute & if it fails attempt fetching of policies from
  // EndPoint attribute
  Priorities  []VolumeReference
}

// VolumeReference refers to references that a volume can use to set its properties,
// policies.
type VolumeReference string

const (
  // ConfigMapVR is the ConfigMap reference type
  ConfigMapVR VolumeReference = "ConfigMap"
  
  // EndPointVR is the EndPoint reference type
  EndPointVR VolumeReference = "EndPoint"
)
```

#### Design -- StorageClass


- NOTE: This is not mandated to be followed at customer deployments
- StorageClass is the placeholder that refers to openebs storage policy & openebs
storage property


```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
   name: sc-cockroachdb-abc:0.5.0
provisioner: openebs.io/provisioner-iscsi
parameters:
    openebs.io/policy: cm-storagepolicy-abc:0.5.0
    openebs.io/property: cm-storageproperty-cstor:0.5.0
```

#### Design -- StoragePolicy -- A K8s ConfigMap

- Maya will use K8s ConfigMap for openebs storage policy specification
  - NOTE: When K8s gets its StoragePolicy Maya will switch to K8s StoragePolicy
- This composes various fine granular storage policies
- Some of the policies can be set to false. A false indicates storage policy 
will not consider that granular policy
- Some of the granular policies are mandatory while others are optional & can be 
set to false

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-storagepolicy-abc:0.5.0
data:
  category: openebs.io/storage-policy
  type: StoragePolicy
  version: 0.5.0
  policies:
    policy.openebs.io/controller-service: cm-ctrl-svc-jiva:0.5.0
    policy.openebs.io/controller-placement: cm-ctrl-plc-jiva:0.5.0
    policy.openebs.io/replica-placement: cm-rep-plc-jiva:0.5.0
    policy.openebs.io/controller-qos: false
    policy.openebs.io/replica-qos: false
    policy.openebs.io/snapshot: false
    policy.openebs.io/cron-jobs: false
    policy.openebs.io/controller-monitoring: false
```

#### Design -- StorageProperty -- A K8s ConfigMap

- Maya will use K8s ConfigMap for openebs storage property
- It composes one (or more) granular properties

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-storageproperty-default:0.5.0
data:
  category: openebs.io/storage-property
  type: StorageProperty
  version: 0.5.0
  properties:
    property.openebs.io/config: cm-cstor-config-default:0.5.0
```

#### Design -- property.openebs.io/config -- A K8s ConfigMap

- Maya will use K8s ConfigMap for property.openebs.io/config
- NOTE: This is a granular storage property

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-config-cstor-default:0.5.0
data:
  category: property.openebs.io/config
  type: CStorConfig
  version: 0.5.0
  yaml: |
    pools:
    - type: src
      create: no
      name: spool
      cachefile: /tmp/cstor/spool.cache
      pooltype:
      diskpaths: /tmp/cstor/sdisk.img
    datasets:
    - type: src
      id: 0 
      create: no
      name: svol
      parent: spool
      volblocksize: 4096
      size: 10737418240
      readonly: off
      sync: always
      copies: 1
      logbias: latency
      compression: on
```


#### Design -- policy.openebs.io/controller-service -- A K8s ConfigMap

- Maya will use K8s ConfigMap for policy.openebs.io/controller-service
- NOTE: There is a granular storage policy

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-ctrl-svc-jiva:0.5.0
data:
  category: policy.openebs.io/controller-service
  type: K8sService
  version: 0.5.0
  yaml: |
    apiVersion: v1
    kind: Service
    metadata:
      name: {{.Name}}-ctrl-svc
    spec:
      ports:
      - name: iscsi
        port: 3260
        protocol: TCP
        targetPort: 3260
      - name: api
        port: 9501
        protocol: TCP
        targetPort: 9501
      selector:
        openebs.io/controller: jiva-controller
        openebs.io/volume: {{.Name}}
      sessionAffinity: None
```

#### Design -- policy.openebs.io/controller-placement -- A K8s ConfigMap

- Maya will use K8s ConfigMap for policy.openebs.io/controller-placement
- NOTE: There is a granular storage policy

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-ctrl-plc-jiva:0.5.0
data:
  category: policy.openebs.io/controller-placement
  type: K8sDeployment
  version: 0.5.0
  yaml: |
    apiVersion: apps/v1beta2 # for versions before K8s 1.8.0 use apps/v1beta1
    kind: Deployment
    metadata:
      name: {{.Name}}-ctrl
      labels:
        openebs.io/controller: jiva-controller
        openebs.io/volume-type: jiva
        openebs.io/volume: {{.Name}}
    spec:
      replicas: 1
      selector:
        matchLabels:
          openebs.io/controller: jiva-controller
          openebs.io/volume: {{.Name}}
      template:
        metadata:
          labels:
            openebs.io/controller: jiva-controller
            openebs.io/volume: {{.Name}}
        spec:
          tolerations:
          - effect: NoExecute
            key: node.alpha.kubernetes.io/notReady
            operator: Exists
            tolerationSeconds: 0
          - effect: NoExecute
            key: node.alpha.kubernetes.io/unreachable
            operator: Exists
            tolerationSeconds: 0
          containers:
          - name: {{.Name}}-ctrl-con
            image: openebs/jiva:0.3-RC2
            args: ["controller", "--frontend", "gotgt", "--clusterIP", "{{.ServiceIP}}", "{{.Name}}"]
            command: ["launch"]
            ports:
            - containerPort: 3260
              protocol: TCP
            - containerPort: 9501
              protocol: TCP
```

#### Design -- policy.openebs.io/replica-placement -- A K8s ConfigMap

- Maya will use K8s ConfigMap for policy.openebs.io/replica-placement
- NOTE: There is a granular storage policy

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-rep-plc-jiva:0.5.0
data:
  category: policy.openebs.io/replica-placement
  type: K8sDeployment
  version: 0.5.0
  yaml: |
    apiVersion: apps/v1beta2 # for versions before K8s 1.8.0 use apps/v1beta1
    kind: Deployment
    metadata:
      name: {{.Name}}-rep
      labels:
        openebs.io/replica: jiva-replica
        openebs.io/volume-type: jiva
        openebs.io/volume: {{.Name}}
    spec:
      replicas: 2
      selector:
        matchLabels:
          openebs.io/replica: jiva-replica
          openebs.io/volume: {{.Name}}
      template:
        metadata:
          labels:
            openebs.io/replica: jiva-replica
            openebs.io/volume: {{.Name}}
        spec:
          tolerations:
          - effect: NoExecute
            key: node.alpha.kubernetes.io/notReady
            operator: Exists
          - effect: NoExecute
            key: node.alpha.kubernetes.io/unreachable
            operator: Exists
          affinity:
            podAntiAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                  matchLabels:
                    openebs.io/replica: jiva-replica
                    openebs.io/volume: {{.Name}}
                topologyKey: kubernetes.io/hostname
          containers:
          - name: {{.Name}}-rep-con
            image: openebs/jiva:0.3-RC2
            args: ["replica", "--frontendIP", "{{.ServiceIP}}", "--size", "{{.Capacity}}", "/openebs"]
            command: ["launch"]
            ports:
            - containerPort: 9502
              protocol: TCP
            - containerPort: 9503
              protocol: TCP
            - containerPort: 9504
              protocol: TCP
```

#### Design -- policy.openebs.io/controller-qos -- A K8s ConfigMap

- TODO

#### Design -- policy.openebs.io/replica-qos -- A K8s ConfigMap

- TODO

#### Design -- policy.openebs.io/snapshot -- A K8s ConfigMap

- TODO

#### Design -- policy.openebs.io/cron-jobs -- A K8s ConfigMap

- TODO

#### Design -- policy.openebs.io/controller-monitoring -- A K8s ConfigMap

- TODO

#### Notes

- Since there will be an explosion of storage policies in form of yamls it makes
sense to have them validated.
- It makes sense to create a CustomResourceDefinition of storage policy that will
have above validation.
  - This can be done after the above implementation
- maya volume controller will take care of provisioning various openebs volumes
- maya volume controller will take care of applying various storage policies to
a volume
- maya volume controller will take care of validating storage policies
- alis-name-type:label is the new naming convention openebs will use to name the
 K8s Kinds that are declared in YAML format
- The namespace will be extracted from ConfigMap's namespace
- maya volume controller may merge some of these storage policies & will control
the ordered execution of these policies
- It also makes sense to create a CustomResourceDefinition of openebs volume
- There will be experiments to implement the openebs volume as a K8s initializer to
a PersistentVolume.
  - In other words, openebs volume CRD object will be created when a PV based on
  openebs volume provisioner gets initialized.
  - This might also mean openebs-provisioner may no more be required
