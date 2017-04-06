## mtool

`mtool` strives for declarative programming to implement a workflow. It 
will not be another DSL library or toolkit. However, it will mention the good 
practices to implement this declarative mode of programming.

The trials of `mtool` concepts will be done in specific folders within `elves` 
repository.

### Expressing intent of the workflow is mtool's key focus

Express the [intent](intent.md) of your workflow in a yaml file.

### Tools to express intent 

- Ansible
- [Sup](../sup/README.md)

## How can one use mtool ?

### As a CLI

- Use a tool that can compose these test cases & expose them as a CLI
  - This CLI can be packaged & released
  - e.g. `mtest` can be a packaged distribution which derives all its design principles from `mtool`.

### Over HTTP

- Expose the CLI over HTTP

### Over SLACK

- SLACK plugin that accepts above CLI commands
