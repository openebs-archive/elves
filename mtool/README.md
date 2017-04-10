## mtool

`mtool` believes in taskrunners that can be declarative . It will not be another 
DSL library or toolkit to achive its beliefs. It will a best practices guideline
to choose one or more declarative tool(s). It will also discuss ways to extend this
declarative taskrunner to satisfy various `DevOps`, `NoOps` related thought processes.

`mtool` makes its possible to practice `continuous delivery` & `continuous improvement` 
in OpenEBS projects.

>NOTE: The trials of `mtool` concepts will be done in specific folders within `elves` 
repository.

### Expressing intent of the taskrunner is mtool's key focus

Express the [intent](intent.md) of your taskflow in a yaml file. 

Some of the benefits:

- `Version control everything`: intents can be version controlled
- `Quality checks`: forces complex logic to be broken down into modules that build the intent
- `Bring the pain forward`: error prone tasks can be debugged during review itself
- `Everyone is responsible`: helps non-programmers to participate in developing intents

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
