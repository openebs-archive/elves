## Integration Testing Best Practices Guide

A contributor needs to understand &/ take care of below items:

### Topics - 1

```yaml
Topics:
  - How to implement integration test cases in the age of containers & kubernetes ?
  - Should you use a tool to achieve this ?
  - How do you monitor the success & failure of your test cases ?
  - Can you use this tool in customer setups (dev, test or prod) ?
```

#### Refer

```yaml
Integration Testing Sample In The Age of Containers & K8s:
  - https://github.com/kubernetes/kubernetes/tree/master/examples/cockroachdb
  - https://github.com/kubernetes/kubernetes/blob/master/examples/cockroachdb/demo.sh:
    - You may implement a simple shell script
    - You might want to write a simple .go file that has main() function
    - You may not want another tool(s) if `kubectl` & host of K8s features suffice
```

