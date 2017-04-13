## Scripting Best Practices Guide

A contributor needs to understand &/ take care of below non-functional features:

### Topics - 1

```yaml
Topics:
  - Why & how to create a clean_up function ?
  - Why & how to create a error_exit function ?
  - Why & how to create safe (race free) temp files ?
  - Why SIGKILL is not the best option for signal handling ?
  - Why quotes are important ?
  - Did you try tracing when you are in trouble ?
```

#### Refer

```yaml
Errors And Signals:
  - http://linuxcommand.org/lc3_wss0150.php
  - Errors And Signals And Traps (Oh, My!) - Part 1
  - Errors And Signals And Traps (Oh, My!) - Part 2
  - Stay Out Of Trouble
```

### Topics - 2

```yaml
Topics:
  - What is linting & why is linting required ?
  - Can we use a tool to lint ?
  - Can the tool be integrated with Travis CI ?
```

#### Refer

```yaml
online lint:
  - https://www.shellcheck.net/
shellcheck:
  - https://github.com/koalaman/shellcheck
lint-shell-scripts:
  - https://carlosbecker.com/posts/lint-shell-scripts/
```

### Topics - 3

```yaml
Topics:
  - What is unit testing ?
  - Do you unit test your scripts ?
```

#### Refer

```yaml
Bash Automated Testing System:
  - https://github.com/sstephenson/bats
```
