## mtool

`mtool` is a concept that encompasses 
[taco bell programming](http://whatis.techtarget.com/definition/Taco-Bell-programming). 
We are just rebranding it as `mtool`.

In future, we might put `mtool` concepts in `openebs` repository. The trials of 
`mtool` concepts will be done in `elves` repository.

## mtool design concepts

### Intent

Express the [intent](intent.md) of your requirement in a yaml file.

### CLI

- Use a tool that can compose these test cases & expose them as a CLI
  - This CLI can be packaged & released
  - e.g. `mtest` can be a packaged distribution which derives all its design principles from `mtool`.

### Over HTTP

- Expose the CLI over HTTP

### Over SLACK

- SLACK plugin that accepts above CLI commands
