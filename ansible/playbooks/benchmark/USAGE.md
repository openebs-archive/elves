# Steps to use the performance benchmark framework 

This file describes the steps to use the benchmark framework. 

## Step-1 Specify the benchmark test variables 

Edit the __Benchmark test-specific params__ section in the benchmark.yml to specify the tool and volume definition. 
Also, set the workload template param to "user"

## Step-2 Specify the workload characteristics

Edit the sample input YAML file at `benchmark/templates/user` with desired params. More than one template can be 
created and placed in this folder. Note that, certain params maybe N/A for the specified tool. For instance, the 
queue depth param is N/A for vdbench. Such params are not parsed by the config generator scripts. Comments are provided 
in the template samples indicating such cases.

Also, it needs to be noted that some of the values specified in the input YAML translates to multiple tool-specific flags. 
For example, with fio, an access pattern of type "Random" causes the flags "norandommap=1" & "randrepeat=0" to be included
in the config file.These ensure offset-ranges are not repeated and I/O history is not referred respectively. This is done to 
simulate real-world behaviour.

## Step-3 Run the benchmark playbook

Run the main benchmark playbook using the command shown below to trigger the benchmark run. 

```
ciuser@OpenEBSClient:~/openebs/e2e/ansible$ ansible-playbook playbooks/benchmark/benchmark.yml 
```

## Step-4 View test results 

The tool log files and plots (IOPS v/s interval, BW v/s interval, Latency v/s interval) are placed in the default log
directory on the host at `/mnt/logs` under the folder created with template name.







