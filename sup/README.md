## sup 

[Sup](https://github.com/pressly/sup) is a tool similar in spirit to Ansible. 
We shall conduct various experiments with `sup`. If these experiments are
satisfactory to the needs of [mtool](../mtool/README.md) concept, we shall 
implement `mtest` using `sup`.

### Hands on with sup

- Title: **Run over network**

```yaml
networks:
    staging:
      hosts:
        - root@myhost
```

- Title: **Run locally**
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

- Title: **Run commands in a single session**
- Command: `sup staging lsroot`
- Expected O/P: display the contents of /

```yaml
commands:
    lsroot:
      desc: list root
      run: cd /; ls
```

- Title: **Run each command in a separate session**
- Command: `sup staging goroot ls`
- Expected O/P: display the contents of $HOME

```yaml
commands:
    goroot:
      desc: go to root path
      run: cd /
    ls:
      desc: list files
      run: ls
```

- Title: **Run multiple commands**
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

- Query: Can we import external sup file(s) ?
  - Answer: This is not possible
- Instead once can follow `bootstrapping approach`

> Bootstrapping is the leveraging of a small initial effort into something larger 
and more significant.

- We shall have a bootstraping Supfile
- A bootstrap Supfile will consist of defaults & constants
- The bootstrap Supfile will consist of a single command only called the bootstrap commad.
- This bootstrap command will:
  - pass all its defaults/constants & 
  - Invoke a sup `target` command present in external supfile
    - NOTE: A sup `target` command is an alias of multiple commands 
  - This process of bootstrapping can be followed again in external sufile(s)


```yaml
 install-omm:
    desc: Installs OpenEBS Maya Master
    local: >
      sup -f ./infra-omm/Supfile $SUP_ENV $SUP_NETWORK install
```
