# Elves

Helpers for setting up the development and demo environements.

This is our `programmable infrastructures` repository.

## Setting up build environement

### Pre-requisites

- Linux Host ( say Ubuntu )
- Virtual Box 
- Vagrant
- Git
- Create a developement folder ( say /opt/dev )

### Steps

In your linux host

```
cd <dev-folder>
sudo git clone https://github.com/openebs/elves.git
cd elves/vagrant
sudo vagrant up
sudo vagrant ssh
```

Once you are in the development VM, you can go to the required project and issue 'make'

```
vagrant@openebs-dev-01:~$ cd /opt/gopath/src/github.com/openebs/openebs/
vagrant@openebs-dev-01:/opt/gopath/src/github.com/openebs/openebs$ sudo make
```
