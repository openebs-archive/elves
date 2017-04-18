## Storage Interface

```yaml
Status:
  - In Progress
Owners:
  - @amitkumardas
Github Repo:
  - NA
```

### Motivation

There has been a requirement to configure, manage & use local as well as remote 
storage as persistent volumes for containers. In some cases, this storage is 
expected to persist beyond the lifecycle of containers. There are also cases, 
where the storage needs to be ephemeral *(i.e. the data only persists during the
lifetime of the container)*.

OpenEBS treats both `LOCAL` as well `REMOTE` storage with equal priority. The
priority is actually determined much above the storage layers.

Local vs. Remote debate might be triggered by:

1. The applications that wants to use these storage e.g. some of the applications 
are inherently distributed and need empheral storage that can utilize the local 
SSDs for performance.

2. The operational aspects (i.e. costs, maintainance, etc) of `networked storage`
 vs. `hyperconverged storage` vs. `containerized storage`.

### Aim of Storage Interface

Storage interface will provide a common layer to containers needing storage. In 
other words this interface will **abstract the local & remote storage features**
and provide a unified approach to expose volumes from the node it is operating on.

Below are some of the items that `Storage Interface` will deal with:

- Will provide data backup via snapshots
- Will provide data recovery by restoring the snapshots
- Will provide data purge/cleanup
- Will provide capacity management
- Will discover storage &/ storage characteristics w.r.t its operational / local node
- Will provide disk management
- Will provide path management (i.e. avoid path collisions). Security ??
- Will export disk, IO related metrics
- Will provide APIs that a storage operator is expected to execute

### What Storage Interface is not ?

- It is not a volume manager &/ controller.
- It will not manage the volume lifecycle.
- It will facilitate volume lifecycle by executing commands received from external 
volume management services.
- It will not understand node affinity.
- It will not provide data availability beyond its operational / local node
