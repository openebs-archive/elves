## Introduction (Work In Progress)

Intent is one of the core mtool design concepts.

## Intent

> The entire intent is specified in yaml file.

Intents can be tailored to suit a particular requirement. Intents are tagged 
against a level (0 or 1). A `level-0` intent signifies the lowest possible
specification which does not have any dependencies. A `level-1` intent signifies it
has some dependencies with other intents. A dependent intent can belong to any 
level.

A `level-0` intent could be:

- Env intent
  - e.g. env-k8s-1.5.5.yaml
  - e.g. env-maya-0.2.yaml
- Env based intents can deal with providing various defaults or common properties.

Similary below are some examples of `level-1` intents:
  
- Infra based intents can deal with deployment, upgrade etc scenarios.
  - infra-k8s-1.5.5.yaml
  - infra-maya-0.2.yaml
- Test based intents can deal with testing above infra intent.
  - test-maya-0.2.yaml

### Packaging these intents as distributable libraries

These intents together can be packaged as a distributable library. For example, 
we can aggregate a few intents and package them as a new library called mtest. 
However, we need to consider other design concepts of mtool.

### A sample infra-maya-0.2 intent specs

- NOTE: This is a level-1 intent which is dependent on: 
  - `env-maya-0.2.yaml`
  - `infra-maya-omm-0.2.yaml`
  - `infra-maya-osh-0.2.yaml`
- Its specs:
  - download omm-artifacts 
  - generate omm-configs
  - install omm
  - uninstall omm
  - verify omm
  - download osh-artifacts 
  - generate osh-configs
  - install osh
  - uninstall osh
  - verify osh

### A sample infra-maya-omm-0.2 intent specs

- NOTE: This is a level-1 intent which is dependent on: 
  - `env-maya-0.2.yaml`
- Its specs:
  - verify downloads
  - verify config
  - verify install
  - verify un-install
