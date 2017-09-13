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

  make init
  make
  make image
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

- Run the image in a Pod with a single instance Deployment:

  ```
  $ kubectl create -f deployment.yaml
  ```

### Clean up

- To stop this test and clean up the pod
  
  ```
  $ kubectl delete deployment tt-jctrl-resch
  ```

#### Deep Insight(s)

> `rest.InClusterConfig()` uses the `Service Account token` mounted inside the 
Pod at `/var/run/secrets/kubernetes.io/serviceaccount` path.

#### Appendix

- OV - OpenEBS Volume
- Jiva - OpenEBS application that does the low level persistent storage operations
- Controller - Jiva application can be a Controller or a Replica
- tt-jctrl-resch - Short form of `Time Taken for Jiva controller rescheduling`
