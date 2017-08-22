# Overview

The primary objective of the performance benchmark test framework is to ensure that openebs storage performance 
remains consistent across releases for various synthetic workloads that simulate standard application behaviour. 
However, it is also expected to aid understanding of the finer I/O performance-related aspects of openebs storage 
and its behaviour upon changes to the multiple influencing factors, thereby enabling users model/tune it appropriately 
for various requirements. Under the hood, the framework is a set of input files, ansible playbooks, python/shell scripts 
and I/O tools such as fio, vdbench.

Upon extension, the framework should also help evaluators to run benchmark tests against the openebs storage volume 
and generate __comparable__ reports, with storage backend, test platform/env, workload type and benchmark tool being 
the variables. 

The sections below explain the important pieces of this framework 

# Workload templates

The workload templates are input YAMLs that describe certain basic and generic I/O workload characteristics. These 
params have been arrived at after perusing the performance/benchamark methodologies and approach taken by the storage 
community while attempting to simulate application behaviour. The table below discusses these params in brief, while
mentioning the vdbench-fio flags for the same


