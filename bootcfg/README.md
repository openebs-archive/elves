### BootCfg Elf

Playground for defining & testing the install/deployment related features 
of openebs maya master (`omm`) & openebs storage host (`osh`). This will make 
use of various config, specification, etc files. Once the intent works out as 
desired, it will be programmed from within maya.

#### Programmable Infrastructure -via- Maya

To cut it short, `Maya` will use 
[bootcfg](https://github.com/coreos/coreos-baremetal) for setting up its 
`master(s)` & `openebs host(s)`. The install/deployment parts of Maya will become
the http / gRPC consumer to bootcfg. One can expect Maya CLI to provide `options`
*i.e. cli arguments* that work on top of `bootcfg specs`. 

### Setting up the dev environment

#### Pre-requisites

> Your laptop can have below software installed:

- Linux Host (e.g. Ubuntu)
- Virtual Box
- Vagrant
- git

#### Steps

- In your laptop
  - cd <some-dev-folder>
  - sudo git clone https://github.com/openebs/elves.git
  - cd elves/bootcfg
  - vagrant up
  - vagrant ssh
- Running bootcfg (*within the vagrant VM*)
- Setting up Maya master
- Setting up Maya OpenEBS host
- Troubleshooting
