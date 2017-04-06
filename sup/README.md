## sup 

[Sup](https://github.com/pressly/sup) is a tool similar in spirit to Ansible. 
We shall conduct various experiments with `sup`. If these experiments are
satisfactory to the needs of [mtool](../mtool/README.md) concept, we shall 
implement `mtest` using `sup`.

### Hands on with sup

- Title: **Settings to run commands over network**

```yaml
networks:
    staging:
      hosts:
        - root@myhost
```

- Title: **Settings to run commands locally**
- Command: ??

```yaml
networks:
  local:
    hosts:
      - localhost

commands:
    prepare:
        desc: Prepare to upload
        local: npm run build
```

- Title: **How to run multiple commands in a single session**
- Command: `sup staging lsroot`
- Expected O/P: display the contents of /

```yaml
commands:
    lsroot:
      desc: list root
      run: cd /; ls
```

- Title: **How to trigger execution of multiple commands in serial ?**
- Command: `sup staging deploy`
- NOTE: 
  - Each command will be run on all hosts in parallel
  - sup will check return status from all hosts, and run subsequent commands on success only 
  - Thus any error on any host will interrupt the process.

```yaml
targets:
    deploy:
        - build
        - pull
        - migrate-db-up
        - stop-rm-run
        - health
        - slack-notify
        - airbrake-notify
```

- Title: **Global environment variables**
- Command: `sup staging echo`

```yaml
# Global environment variables
env:
  NAME: api
  IMAGE: example/api

commands:
  echo:
    desc: Print some env vars
    run: echo $NAME $IMAGE
```

- Title: **Sup's default environment variables**
- Command: `sup staging echo`

```yaml
commands:
  echo:
    desc: Print some env vars
    run: echo $SUP_NETWORK $SUP_HOST $SUP_ENV
```


## Structuring projects with sup

We can follow `bootstrapping` approach to structure large projects with sup.
This approach assumes significance, since Sup does not allow import of external
Supfiles.

> Bootstrapping is the leveraging of a small initial effort into something larger 
and more significant.

How to employ `bootstrapping` approach ?

- We shall have a Supfile referred to as a bootstrap Supfile.
- A bootstrap Supfile will consist of defaults & constants & other common properties.
- The bootstrap Supfile will consist of a single command referred to as the bootstrap command.
- This bootstrap command will:
  - Pass all its defaults, constants, etc properties to external Supfile(s).
  - Invoke a sup command present in external Supfile:
    - This sup command available in external Supfile is referred to as `target`.
    - A sup `target` command is an alias of multiple commands.
    - i.e. can invoke multiple commands in succession.
- This process of bootstrapping can be followed again in external Supfile(s) as well.

Below is a sample bootstrap Supfile

```yaml
 install-omm:
    desc: Installs OpenEBS Maya Master
    local: >
      sup -f ./infra-omm/Supfile $SUP_ENV $SUP_NETWORK install
```
