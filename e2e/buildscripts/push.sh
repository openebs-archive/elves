#!/bin/bash
set -e

# Get the numeric ID for this docker image tagged with ci
# 
# NOTE:
#   Check the GNUMakefile of this project which is responsible 
# to generate this image in the first place
IMAGEID=$( sudo docker images -q openebs/m-e2e:ci )

if [ ! -z "${DNAME}" ] && [ ! -z "${DPASS}" ]; 
then 
  sudo docker login -u "${DNAME}" -p "${DPASS}"; 
  # Push the development build image i.e. the one tagged with ci
  # to docker hub repository. This is done by default
  sudo docker push openebs/m-e2e:ci ;

  # NOTE:
  #   When github is tagged with a release, Travis will 
  # hold the release tag in env TRAVIS_TAG
  if [ ! -z "${TRAVIS_TAG}" ] ; 
  then
    # Tag this ci image with github release & push it
    sudo docker tag ${IMAGEID} openebs/m-e2e:${TRAVIS_TAG}
    sudo docker push openebs/m-e2e:${TRAVIS_TAG}; 
    # Tag this ci image with latest & push it
    sudo docker tag ${IMAGEID} openebs/m-e2e:latest
    sudo docker push openebs/m-e2e:latest; 
  fi;
else
  echo "WARN: No docker credentials provided.";
  echo "WARN: Skipping uploading openebs/m-e2e:ci to docker hub.";
  echo -e "NOTE: If this was a manual attempt:";
  echo -e "\t1. Login to docker hub.";
  echo -e "\t2. Repeat this action from the same terminal session.";
fi;
