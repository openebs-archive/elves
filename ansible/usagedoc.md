## Role: [k8s-master](k8s-master)
> **Author:** yudaykiran

> **Description:** This role performs the setup the k8s master on the specified host

> **Task:** Get Current User

> **Task:** Create Directory

> **Task:** Change Folder Permissions

> **Task:** Copy Pod YAML to remote

> **Task:** Copy Script to remote

> **Task:** kubeadm reset before init

> **Task:** kubeadm init

> **Task:** Get Current User

> **Task:** Copy k8s credentials to $HOME

> **Task:** Change file permissions

> **Task:** Update KUBECONFIG in .profile

> **Task:** Patch kube-proxy for CNI Networks


## Role: [k8s-localhost](k8s-localhost)
> **Author:** ksatchit

> **Description:** This role prepares the Ansible host machine to undertake execution of k8s roles

> **Task:** Install APT Packages

> **Task:** Install PIP Packages

> **Task:** Create files directory in kubernetes role

> **Task:** Get .deb kubernetes packages google cloud

> **Task:** Check whether deb packages and images are already downloaded

> **Task:** Fetch configure_K8s scripts into k8s-master role

> **Task:** Change Kubernetes version in configure scripts

> **Task:** Fetch the weave .yaml template from GitHub

> **Task:** Fetch configure_K8s scripts into k8s-host role

> **Task:** Get Current User

> **Task:** Change Folder Permissions


## Role: [kubernetes](kubernetes)
> **Author:** yudaykiran

> **Description:** This role prepares the environment for k8s setup on specified hosts

> **Task:** Get Package Updates

> **Task:** Add apt-key

> **Task:** Add Kubernetes apt repository

> **Task:** Get Package Updates

> **Task:** Install APT Packages

> **Task:** Get Current User

> **Task:** Create Directory

> **Task:** Change Folder Permissions

> **Task:** Copy TAR to remote

> **Task:** Copy local deb files to K8s-Master and K8s-Minions


## Role: [k8s-openebs-operator](k8s-openebs-operator)
> **Author:** ksatchit

> **Description:** This role installs the openebs operator (maya-apiserver & openebs-provisioner) on the k8s cluster 

> **Task:** Download YAML for openebs operator

> **Task:** Get kubernetes master name

> **Task:** Get kubernetes master status

> **Task:** None

> **Task:** Ending Playbook Run - K8S master is NOT READY

> **Task:** Deploy the openebs operator yml

> **Task:** Confirm maya-apiserver pod is running

> **Task:** Create fact for pod name

> **Task:** Confirm that maya-cli is available in the maya-apiserver pod

> **Task:** Download YAML for openebs storage classes

> **Task:** Deploy the openebs storageclasses yml

> **Task:** Confirm that openebs storage classes are created


## Role: [master](master)
> **Author:** yudaykiran

> **Description:** This role performs the setup of Maya master on specified host

> **Task:** Setup Maya

> **Task:** Update INI file

> **Task:** Update INI file


## Role: [fio](fio)
> **Author:** ksatchit

> **Description:** This role attaches the openebs storage as a volume to an FIO  application container 

> **Task:** Establish ISCSI Connection

> **Task:** Identify Block Device

> **Task:** Create a File System

> **Task:** Mount Device by Label

> **Task:** fio docker image pull

> **Task:** fio docker instantiate


## Role: [k8s-openebs-cleanup](k8s-openebs-cleanup)
> **Author:** ksatchit

> **Description:** This role performs a cleanup of the openebs operator from the k8s cluster

> **Task:** Delete the openebs storage classes

> **Task:** Confirm storage classes have been deleted

> **Task:** Delete the openebs operator

> **Task:** Confirm pod has been deleted


## Role: [k8s-hosts](k8s-hosts)
> **Author:** yudaykiran

> **Description:** This roles performs a setup of the k8s hosts and adds them into the k8s cluster 

> **Task:** Get Master Status

> **Task:** None

> **Task:** Ending Playbook Run - Master is not UP.

> **Task:** Get Token Name from Master

> **Task:** Get Token from Master

> **Task:** Get Cluster IP from Master

> **Task:** Save Token and Cluster IP

> **Task:** Get the openebs-iscsi flexvolume driver

> **Task:** Create flexvol driver destination

> **Task:** Copy script to flexvol driver location

> **Task:** Reset kubeadm

> **Task:** Setup k8s Hosts


## Role: [iometer](iometer)
> **Author:** ksatchit

> **Description:** This role attaches the openebs disk as a device to an iometer (dynamo) application container

> **Task:** Establish ISCSI Connection

> **Task:** Identify Block Device

> **Task:** IOmeter docker image build

> **Task:** Get localhost's public IP address

> **Task:** IOmeter docker instantiate


## Role: [common](common)
> **Author:** yudaykiran

> **Description:** This role prepares the hosts for installation of OpenEBS storage services 

> **Task:** apt-get update packages

> **Task:** Install APT Packages

> **Task:** Create Directory

> **Task:** Download Maya

> **Task:** Unzip Maya


## Role: [inventory](inventory)
> **Author:** yudaykiran, ksatchit

> **Description:** This role generates the Ansible inventory ('hosts') file and performs a dynamic refresh

> **Task:** Generate Inventory

> **Task:** Sync Inventory


## Role: [volume](volume)
> **Author:** ksatchit

> **Description:** This role creates the volume container on openebs storage hosts 

> **Task:** Create Volume

> **Task:** Get Volume Info

> **Task:** Fetch Target Portal

> **Task:** Fetch Target IQN

> **Task:** Update Target Details


## Role: [hosts](hosts)
> **Author:** yudaykiran

> **Description:** This role performs a setup of the openebs storage hosts and adds them to the cluster 

> **Task:** Get Master Status

> **Task:** None

> **Task:** Ending Playbook Run - Master is not UP.

> **Task:** Setup Maya


## Role: [cleanup](cleanup)
> **Author:** ksatchit

> **Description:** This role tears down the iSCSI sessions on the client & stops the volume container 

> **Task:** Wait {{ io_run_duration }} sec for I/O completion

> **Task:** Unmount the ext4 filesystem

> **Task:** Tear down iSCSI sessions

> **Task:** Remove stale node entries for ISCSI target

> **Task:** Tear down the storage containers


## Role: [vagrant](vagrant)
> **Author:** ksatchit

> **Description:** This role performs environment setup on Vagrant VMs needed to run Openebs storage services

> **Task:** Update Nomad IP in .profile

> **Task:** Update M-APIServer IP in .profile


## Role: [localhost](localhost)
> **Author:** ksatchit

> **Description:** This role prepares the openebs client to run appliation containers 

> **Task:** Install APT Packages

> **Task:** Install PIP Packages


## Role: [prerequisites](prerequisites)
> **Author:** ksatchit

> **Description:** This node prepares the Ansible host before execution of the inventory role

> **Task:** Install PIP Packages


## Role: [ara](ara)
> **Author:** ksatchit

> **Description:** This role sets up the Ansible Run Analysis (ARA) to record playbook runs

> **Task:** Install ara apt packages

> **Task:** Install ara pip packages

> **Task:** Get ara location

> **Task:** Copy ara callback plugin to openebs custom plugins location

> **Task:** Copy ara action plugins to openebs custom plugins location

> **Task:** Copy ara modules to openebs custom plugins location

> **Task:** Update ansible.cfg with ara module library

> **Task:** Update ansible.cfg with ara action plugins

> **Task:** Update ansible.cfg with ara notification callback

> **Task:** Get ara-manage binary location

> **Task:** Start ara webserver on localhost

> **Task:** Display ara UI URL



Generated by [ansible-docgen](https://www.github.com/starboarder2001/ansible-docgen)
