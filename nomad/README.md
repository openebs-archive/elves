# NOMAD Elf

Playground for defining and testing nomad jiva job specs.

## Setting up build/development environement

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
cd elves/nomad
sudo vagrant up
sudo vagrant ssh
```

Once you are in the development VM, you can run your nomad as follows:

```
sudo nomad agent -dev -config /vagrant/config
```
