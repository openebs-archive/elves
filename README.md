# Elves

Helpers for setting up the development and demo environements.

This is our `programmable infrastructures` repository.

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
cd elves/dev
sudo vagrant up
sudo vagrant ssh
```

Once you are in the development VM, you can go to the required project and issue 'make'

```
vagrant@openebs-dev-01:~$ cd /opt/gopath/src/github.com/openebs/longhorn/
vagrant@openebs-dev-01:/opt/gopath/src/github.com/openebs/longhorn$ sudo make
```

## FAQ

### When I do "vagrant up", I get an error "Insecure world writable dir"
Check the permissions of the vagrant directory. Usually located in your "~/.vagrant.d/". You will notice this error when you interchange between running vagrant as sudo and non-sudo users. You can set the permissions and ownership of the vagrant cache dir "~/.vagrant.d/" to your user name.
```
sudo chown -R <user>:<user> ~/.vagrant.d/
```

### Unable to install the vagrant caching plugin. 
Make sure the vagrant version is 1.8.4 and above. On ubuntu, you might have to download the latest version (..deb file) from the vagrant releases page (https://www.vagrantup.com/downloads.html) and click on it. 


To speed up the subsequent launching of the vagrant vms, caching plugin needs to be installed. 
