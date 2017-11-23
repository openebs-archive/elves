# Setting up gcloud-sdk locally

This document provides instructions for setting up gcloud-sdk and kubectl on a local workstation.

## Creating a Service Account.(Admin Only)

This section provides instructions to set up **Service Accounts** to access your clusters.

1. Start your browser.
2. Open _Console_ on __Google Cloud Platform__.
3. Under _IAM & admin_, select __Service accounts__.
4. Click _Create Service Account_.
5. In the _Create service account_ window:
  - For __Service account name__, enter a meaningful name.(Ex: Developer, Admin, ...)
  - For __Role__ select _Container->Container Engine Developer_.
  - Tick the __Furnish a new private key__ checkbox.
  - Click __Create__ and save the private key file on your local workstation.

## Setting up environment for Vagrant VM

1. Use a Vagrantfile that uses ubuntu/xenial and bring up the Vagrant VM.
2. Copy the private key to the HOME directory of the Vagrant VM.
3. Open the .profile file and export the following:
    
    ```bash
    export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)"
    export KEY_FILE=<private_key_file_path> #Example: /home/ubuntu/example_key.json
    export CLOUDSDK_CORE_PROJECT=<gcloud project id> #Example: strong-thor-764539
    ```

4. Run the following command to save the settings:

    ```bash
    source ~/.profile
    ```

5. Run the following command to create a script file:

    ```bash
    vi setup_gcloud_sdk.sh
    ```

6. Copy the following script to the script file created above and save it.

    ```bash
    #!/bin/bash

    IS_GCLOUD_INSTALLED=$(which gcloud>> /dev/null 2>&1; echo $?)
    if [ $IS_GCLOUD_INSTALLED -eq 0 ]; then
         echo "gcloud sdk is installed; Skipping"
         sleep 2
    else

         echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
         curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
         sudo apt-get update && sudo apt-get install -y google-cloud-sdk kubectl
         gcloud auth activate-service-account --key-file=$KEY_FILE
    fi
    ```

7. Change the permission of the script file to executable.

    ```bash
    chmod +x setup_gcloud_sdk.sh
    ```

8. Run the script.

    ```bash
    ./setup_gcloud_sdk.sh
    ```

## Setting up kubeconfig for your cluster


1. Run the following command to get the list of clusters:

    ```bash
    gcloud container clusters list
    ```

2. Run the following command to create a kubeconfig for cluster __Example-Cluster__

    ```bash
    gcloud container clusters get-credentials Example-Cluster --zone us-east1-a
    ```

    _Note: The zone should match with the zone, where the cluster was created._

3. Run __kubectl__ commands to verify.

    ```bash
    kubectl get nodes
    ```
