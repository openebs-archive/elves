## mtool

`mtool` believes in taskrunners that can be declarative . It will not be another 
DSL library or toolkit to achive its beliefs. It will a best practices guideline
to choose one such declarative tool. It will also discuss ways to extend this
declarative taskrunner to satisfy various `devops` thought processes.

The trials of `mtool` concepts will be done in specific folders within `elves` 
repository.

### Expressing intent of the taskrunner is mtool's key focus

Express the [intent](intent.md) of your taskflow in a yaml file.

### Tools to express intent

- Ansible
- [Sup](../sup/README.md)
- [Alfred](https://github.com/kcmerrill/alfred/)

## Extending the reach of this declarative taskrunner ?

### As a CLI

- Use a tool that can compose these tasks & expose them as a CLI
  - This CLI can be packaged & released
  - e.g. `mtest` can be a packaged distribution which derives all its design principles from `mtool`.

### Over HTTP

- Expose the CLI over HTTP

### Over SLACK

- SLACK plugin that accepts above CLI commands
