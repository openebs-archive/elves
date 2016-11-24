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

## Initial Setup

In your linux host
```
cd <dev-folder>
sudo git clone https://github.com/openebs/elves.git
cd elves/nomad
vagrant up
vagrant ssh
```

## Running nomad
```
cd <dev-folder>/elves/nomad
vagrant up
vagrant ssh
vagrant@nomad-dev:~$ nomad agent -dev -config /vagrant/config
```

## Running jiva 
```
vagrant@nomad-dev:~$ nomad status
vagrant@nomad-dev:~$ nomad run /vagrant/jobs/simple-jiva-vol.nomad 
vagrant@nomad-dev:~$ docker images
```

## Checking status of trouble shooting
```
vagrant@nomad-dev:~$ nomad status
vagrant@nomad-dev:~$ docker images
vagrant@nomad-dev:~$ docker ps -a
```

Check docker logs
```
sudo journalctl -fu docker.service
```


