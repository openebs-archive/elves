# Overview

The primary objective of the performance benchmark test framework is to ensure that OpenEBS storage performance 
remains consistent across releases for various synthetic workloads that simulate standard application behaviour. 
However, it is also expected to aid understanding of the finer I/O performance-related aspects of OpenEBS storage 
and its behaviour upon changes to the multiple influencing factors, thereby enabling users model/tune it appropriately 
for various requirements. Under the hood, the framework is a set of input files, ansible playbooks, python/shell scripts 
and I/O tools such as fio, vdbench.

It can currently filesystem?? benchmark tests on OpenEBS storage via fio or vdbench tools. Upon extension, the 
framework should also help evaluators to run benchmark tests against the OpenEBS storage volume and generate __comparable__ 
reports, with storage backend, test platform/env, workload type and benchmark tool being the variables. 

The following sections explain the important pieces of this framework and the benchmark directory structure.

# Workload templates

The workload templates are input YAMLs that describe basic and advanced I/O workload characteristics. These 
params are arrived at after pursuing the performance/benchamark methodologies and approach taken by the storage 
community while attempting to simulate application behavior. The following table discusses some of the fs related params 
in brief, while mentioning the vdbench-fio flags for the same.

| Workload Param       |   fio flag                  | vdbench flag        |       Notes                                               |
| --------------       | ----------------------------| ------------        | ----------------------------------------------------------|
|  R/W ratio           | readmixread, readmixwrite   | readpct             | Percentage mix of reads, writes                           |
|  I/O size            | blocksize                   | xfersize            | Size of each I/O block                                    | 
|  Access pattern      | fileio                      | readwrite           | Type of I/O access, i.e., random, sequential              | 
|  Worker count        | threads                     | numjobs             | No of threads/processes performing the same workload      | 
|  Queue depth         | n/a (threads)               | iodepth             | No of outstanding/queued I/Os                             |
|  I/O synchronization | ioengine, sync              | openflags           | How I/O is issued by the job, i.e., as a sync or async    |
|  I/O buffering       | direct, buffered            | n/a (openflags)     | Determines whether host fs/kernel page cache is used      | 
|  Data reduction      | dedupe %, buffer_compress_% | dedupratio,compratio| Reducible data sent as part of the I/O                    | 
|  I/O arrival         | rate_process_               | distribution        | Specifies I/O arrival rate, burst I/O options             | 

The workload templates are placed in the `benchmark/templates` directory under the following sub-folders: 
- `Basic` (simple test I/O profiles), 
- `Standard` (application simulation workloads, run as part of CI)  
- `User` (placeholder for user-defined/custom workloads) for potential storage evaluators

# Config generators 

The generic workload templates are converted into tool specific config files by python scripts fio-confgen.py and vdbench-confgen.py,
placed in `benchmark/fio` and `benchmark/vdbench` folders respectively. These scripts also inject some additional/support flags apart 
from the ones described in the table above, based on the param values. 

# Test variables input file 

The benchmark test variables, such as the benchmark tool, workload template type, environment aspects (mount points, log directories), 
supporting files, links, user credentials etc., can be configured in this input YAML file - `benchmark/benchmark-vars.yml` 

# Test containers 

The tool-specific config files generated are run against the storage through test containers which run the specified tool, i.e., fio, 
vdbench etc. The run logs and respective plots are placed in the specified bind mounted host directory. The Dockerfile and 
supporting directories/files are placed at `benchmark/fio` and `benchmark/vdbench`

# Playbooks 

A master playbook `benchmark/benchmark.yml` triggers the benchmark tests by first installing the prerequisites on the storage hosts (`benchmark-prerequisites.yml`) followed by the following:

- provisioning the storage volume (`common/benchmark-provision.yml`) 
- generating tool-specific config files and running the test containers (`fio/benchmark-fio.yml` or `vdbench/benchmark-vdbench.yml`)  
- cleaning up the environment (`benchmark-cleanup.yml`) 









