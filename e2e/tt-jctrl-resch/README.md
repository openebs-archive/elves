# Time Taken for jiva controller rescheduling

This test logic calculates the time taken for OV Jiva controller during
re-scheduling.

## Building this e2e logic

First compile the application for Linux:

  ```
  cd tt-jctrl-resch
  GOOS=linux go build -o ./app .
  ```
    
Then package it to a docker image using the provided Dockerfile to run it on
Kubernetes.

If you are running a Docker engine, you can build this image directly.

  ```
  docker build -t openebs/m-e2e-tt-jctrl-resch:ci .
  ```

## Alternative approach to build & image

Make use of the Vagrantfile available at ../e2e

  ```bash
  vagrant up
  vagrant ssh
  cd tt-jctrl-resch

  make
  ```

## World of Kubernetes

Push it to a registry that your Kubernetes cluster can pull from.

e.g. https://hub.docker.com/r/openebs/m-e2e-tt-jctrl-resch/tags/

#### Pre-requisite

- Run OpenEBS operator yaml

  ```
  $ kubectl create -f openebs-operator.yaml
  ```

#### Run

- Run this e2e image as a Kubernetes Job

  ```
  $ kubectl create -f deployment.yaml
  ```

### Clean up

- To delete this job
  
  ```
  $ kubectl delete job e2e-tt-jctrl-resch
  ```

- Check if any orphaned Pod(s) exist

  ```
  $ kubectl get pod -a
  ```

#### Appendix

- OV - OpenEBS Volume
- Jiva - OpenEBS application that does the low level persistent storage operations
- Controller - Jiva application can be a Controller or a Replica
- tt-jctrl-resch - Short form of `Time Taken for Jiva controller rescheduling`
