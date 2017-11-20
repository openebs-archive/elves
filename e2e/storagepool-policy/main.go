/*
Copyright 2017 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// FetchCommand will fetch a specific property value from a
// K8s Kind
type FetchCommand string

const (
	// VolumeFromPVCFC fetches volume name i.e. pv from the pvc
	VolumeFromPVCFC FetchCommand = `kubectl get pvc demo-vol1-claim -o jsonpath='{.spec.volumeName}' --namespace=default`

	// MayaServiceIPFC fetches maya api service's IP address
	MayaServiceIPFC FetchCommand = `kubectl get svc maya-apiserver-service -o jsonpath='{.spec.clusterIP}'`
)

// VerifyCommand represents a command used for verifying availability
// of K8s Kinds
type VerifyCommand string

const (
	// MayaServiceAccountVC verifies presence of maya service account
	MayaServiceAccountVC VerifyCommand = `kubectl get sa openebs-maya-operator -o name`

	// MayaClusterRoleVC verifies presence of maya cluster role
	MayaClusterRoleVC VerifyCommand = `kubectl get clusterrole openebs-maya-operator -o name`

	// MayaClusterRoleBindingVC verifies presence of maya cluster role binding
	MayaClusterRoleBindingVC VerifyCommand = `kubectl get clusterrolebinding openebs-maya-operator -o name`

	// PerconaStorageClassVC verifies presence of percona storage class
	PerconaStorageClassVC VerifyCommand = `kubectl get sc openebs-percona -o name`

	// MayaAPIServiceVC verifies presence of maya api service
	MayaAPIServiceVC VerifyCommand = `kubectl get svc maya-apiserver-service -o name`

	// OpenEBSProvisionerVC verifies presence of openebs provisioner
	OpenEBSProvisionerVC VerifyCommand = `kubectl get deploy openebs-provisioner -o name`

	// HostDirStoragePoolVC verifies presence of host dir storage pool
	HostDirStoragePoolVC VerifyCommand = `kubectl get sp sp-hostdir -o name`
)

// IsRunningCommand represents a command used to verify if K8s Kind
// is in `Running` status
type IsRunningCommand string

const (
	// MayaAPIServiceIRC verifies if maya api service is running
	MayaAPIServiceIRC IsRunningCommand = `kubectl get po -l name=maya-apiserver -o jsonpath='{.items[0].status.phase}'`

	// OpenEBSProvisionerIRC verifies if openebs provisioner is running
	OpenEBSProvisionerIRC IsRunningCommand = `kubectl get po -l name=openebs-provisioner -o jsonpath='{.items[0].status.phase}'`

	// PerconaAppIRC verifies if percona app is running
	PerconaAppIRC IsRunningCommand = `kubectl get po -l name=percona -o jsonpath='{.items[0].status.phase}'`
)

// DisplayCommand represents a command used to display K8s Kinds
type DisplayCommand string

const (
	// PerconaStorageClassDC will display the details of percona
	// storage class
	PerconaStorageClassDC DisplayCommand = `kubectl get sc openebs-percona -o yaml`

	// OpenEBSVolumeRepDC will display the details of OpenEBS volume replica(s)
	OpenEBSVolumeRepDC DisplayCommand = `kubectl get deploy -l openebs/replica=jiva-replica`

	// OpenEBSVolumeCtrlDC will display the details of OpenEBS volume controller(s)
	OpenEBSVolumeCtrlDC DisplayCommand = `kubectl get deploy -l openebs/controller=jiva-controller`

	// PerconaPVCDC will display the details of percona pvc
	PerconaPVCDC DisplayCommand = `kubectl get pvc demo-vol1-claim -o yaml`
)

// DeleteCommand represents a command used to delete K8s Kinds
type DeleteCommand string

const (
	// PerconaStorageClassDELC will delete the percona storage class
	PerconaStorageClassDELC DeleteCommand = `kubectl delete sc openebs-percona`

	// PerconaAppDELC will delete the percona pod
	PerconaAppDELC DeleteCommand = `kubectl delete po percona`

	// PerconaPVCDELC will delete the percona pvc
	PerconaPVCDELC DeleteCommand = `kubectl delete pvc demo-vol1-claim`

	// PerconaSCDELC will delete the percona sc
	PerconaSCDELC DeleteCommand = `kubectl delete sc openebs-percona`

	// HostDirSPDELC will delete the host dir storage pool
	HostDirSPDELC DeleteCommand = `kubectl delete sp sp-hostdir`

	// StoragePoolCRDDELC will delete the storage pool crd definition
	StoragePoolCRDDELC DeleteCommand = `kubectl delete crd storagepools.openebs.io`
)

// PVDELC will delete the provided PV
var PVDELC = `kubectl delete pv %s`

// PerconaOVDELC will delete the provided openebs volume
var PerconaOVDELC = `curl http://%s:5656/latest/volumes/delete/%s`

// CreateCommand represents a command used to create K8s Kinds
type CreateCommand string

const (
	// PerconaStorageClassCC when run will create a K8s storage class
	// with replica count = 2
	PerconaStorageClassCC CreateCommand = `cat <<EOF | kubectl create -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: openebs-percona
provisioner: openebs.io/provisioner-iscsi
parameters:
  pool: hostdir-var
  openebs.io/jiva-replica-count: "1"
  openebs.io/capacity: "2G"
  openebs.io/jiva-replica-image: "openebs/jiva:0.4.0"
  openebs.io/storage-pool: "sp-hostdir"
EOF`

	// StoragePoolCRDCC when run will create a K8s StoragePool
	// CRD definition
	StoragePoolCRDCC CreateCommand = `cat <<EOF | kubectl create -f -
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: storagepools.openebs.io
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: openebs.io
  # version name to use for REST API: /apis/<group>/<version>
  version: v1alpha1
  # either Namespaced or Cluster
  scope: Cluster
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: storagepools
    # singular name to be used as an alias on the CLI and for display
    singular: storagepool
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: StoragePool
    # shortNames allow shorter string to match your resource on the CLI
    shortNames:
    - sp
EOF`

	// HostDirStoragePoolCC when run will create a K8s StoragePool
	HostDirStoragePoolCC CreateCommand = `cat <<EOF | kubectl create -f -
apiVersion: openebs.io/v1alpha1
kind: StoragePool
metadata:
  name: sp-hostdir
  type: hostdir
spec:
  path: "/var/coolebs" 
EOF`

	// MayaAPIServiceCC when run will create maya api service
	MayaAPIServiceCC CreateCommand = `cat <<EOF | kubectl create -f -
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: maya-apiserver
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: maya-apiserver
    spec:
      serviceAccountName: openebs-maya-operator
      containers:
      - name: maya-apiserver
        imagePullPolicy: Always
        image: openebs/m-apiserver:crd
        ports:
        - containerPort: 5656
---
apiVersion: v1
kind: Service
metadata:
  name: maya-apiserver-service
spec:
  ports:
  - name: api
    port: 5656
    protocol: TCP
    targetPort: 5656
  selector:
    name: maya-apiserver
  sessionAffinity: None
EOF`

	// OpenEBSProvisionerCC when run will create openebs provisioner
	OpenEBSProvisionerCC CreateCommand = `cat <<EOF | kubectl create -f -
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: openebs-provisioner
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: openebs-provisioner
    spec:
      serviceAccountName: openebs-maya-operator
      containers:
      - name: openebs-provisioner
        imagePullPolicy: Always
        image: satyamz/provisioner:v0.5-rc
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
EOF`

	// PerconaAppCC when run will create a percona app
	PerconaAppCC = `cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Pod
metadata:
  name: percona
  labels:
    name: percona
spec:
  containers:
  - resources:
      limits:
        cpu: 0.5
    name: percona
    image: percona
    args:
      - "--ignore-db-dir"
      - "lost+found"
    env:
      - name: MYSQL_ROOT_PASSWORD
        value: k8sDem0
    ports:
      - containerPort: 3306
        name: percona
    volumeMounts:
    - mountPath: /var/lib/mysql
      name: demo-vol1
  volumes:
  - name: demo-vol1
    persistentVolumeClaim:
      claimName: demo-vol1-claim
---
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: demo-vol1-claim
spec:
  storageClassName: openebs-percona
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5G
EOF`
)

// execSh executes the provided command in sh
//
// NOTE: This function has been adapted from
// https://github.com/c-bata/kube-prompt/blob/master/kube/executor.go
func execSh(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		fmt.Printf("[WARN] Missing command\n")
		return
	}

	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[ERR] %s\n", err.Error())
	}
}

// prompt prompts for user input. Here typing any command
// will continue the execution. This works only
// when this file is run as `go run main.go`
//
// NOTE:
//    This is very helpful for demo purposes. One can inject
// prompt() calls in between various invocations and explain
// what is happening or what is expected to happen during these
// user prompts.
func prompt(msg string) {
	fmt.Printf("[PROMPT] " + msg + " ?")

	var i int
	fmt.Scanf("%d\n", &i)
}

// execShResult will run the command in sh & return its result
//
// NOTE: This function has been adapted from
// https://github.com/c-bata/kube-prompt/blob/master/kube/executor.go
func execShResult(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("Missing command")
	}

	out := &bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	r := string(out.Bytes())
	return r, nil
}

//
func display(header string, dc DisplayCommand) error {
	fmt.Printf("\n[INFO] " + header + "\n")

	op, err := execShResult(string(dc))
	if err != nil {
		return err
	}

	fmt.Printf("\n---\n")
	fmt.Printf("\n%s\n", op)
	fmt.Printf("\n---\n")

	return nil
}

//
func displayAll(header string, dcs ...DisplayCommand) error {
	for _, dc := range dcs {
		err := display(header, dc)
		if err != nil {
			return err
		}
	}

	return nil
}

// verifyPresence verifies the presence of a K8s Kind object
func verifyPresence(vc VerifyCommand) bool {
	op, err := execShResult(string(vc))
	if err != nil {
		fmt.Printf("[WARN] %s\n", err.Error())
		return false
	}

	if len(op) == 0 {
		fmt.Printf("[WARN] Nil output for command '%s'\n", string(vc))
		return false
	}

	return true
}

// fetch gets the value of a specific property from a K8s
// Kind object
func fetch(fc FetchCommand) (string, error) {
	op, err := execShResult(string(fc))
	if err != nil {
		return "", err
	}

	return op, nil
}

// verifyRunning verifies if a specifc K8s Kind object is
// in a running state
func verifyRunning(irc IsRunningCommand) bool {

	// 120*3 seconds timeout
	repeat := 120
	delay := 3 * 1000 * time.Millisecond
	for i := 1; i <= repeat; i++ {
		op, err := execShResult(string(irc))
		if err != nil {
			fmt.Printf("[WARN] %s\n", err.Error())
			return false
		}

		if op == "Running" {
			return true
		}

		// Sleep a bit
		time.Sleep(delay)
	}

	return false
}

//
func verifyAllPresence(vcs ...VerifyCommand) bool {
	for _, vc := range vcs {
		if !verifyPresence(vc) {
			return false
		}
	}
	return true
}

//
func create(cc CreateCommand) {
	execSh(string(cc))
}

//
func delete(dc DeleteCommand) {
	execSh(string(dc))
}

type envkey string

const (
	// DemoModeEK is the ENV variable that flags if this e2e should
	// run as a demo
	DemoModeEK envkey = "E2E_IS_DEMO"
)

// truthyValues maps a set of values which are considered as true
var truthyValues = map[string]bool{
	"1":    true,
	"YES":  true,
	"TRUE": true,
	"OK":   true,
}

// CheckTruthy checks for truthiness of the passed argument.
func CheckTruthy(truth string) bool {
	return truthyValues[strings.ToUpper(truth)]
}

// getEnv fetches the environment variable value from the machine's
// environment
func getEnv(ek envkey) string {
	return strings.TrimSpace(os.Getenv(string(ek)))
}

// main has the entire business logic
func main() {
	// PRE-REQUISITES -- These cannot be automated in this e2e logic !!
	// verify maya service account
	if !verifyPresence(MayaServiceAccountVC) {
		fmt.Printf("[ERR] Missing maya service account\n")
		return
	}

	if !verifyPresence(MayaClusterRoleVC) {
		fmt.Printf("[ERR] Missing maya cluster role\n")
		return
	}

	if !verifyPresence(MayaClusterRoleBindingVC) {
		fmt.Printf("[ERR] Missing maya cluster role binding\n")
		return
	}

	// SETUP -- Even if not available, e2e logic has handled this !!
	if !verifyPresence(MayaAPIServiceVC) {
		fmt.Printf("[INFO] Missing maya api service. Will create one.\n")
		create(MayaAPIServiceCC)
	}
	if !verifyPresence(OpenEBSProvisionerVC) {
		fmt.Printf("[INFO] Missing openebs provisioner. Will create one.\n")
		create(OpenEBSProvisionerCC)
	}

	// verify if maya api service is running
	isMayaAPIRunning := verifyRunning(MayaAPIServiceIRC)
	if !isMayaAPIRunning {
		fmt.Printf("[ERR] Maya api service is not running\n")
		return
	}

	// verify if openebs provisioner is running
	isOPRunning := verifyRunning(OpenEBSProvisionerIRC)
	if !isOPRunning {
		fmt.Printf("[ERR] OpenEBS provisioner is not running\n")
		return
	}

	// this E2E instance's LOGIC starts from here
	prompt("Create host dir storage pool CRD & its impl")
	create(StoragePoolCRDCC)
	create(HostDirStoragePoolCC)
	if !verifyPresence(HostDirStoragePoolVC) {
		fmt.Printf("[ERR] Host dir storage pool could not be created\n")
		return
	}

	prompt("Create percona storage class")
	create(PerconaStorageClassCC)
	// verify percona storage class
	if !verifyPresence(PerconaStorageClassVC) {
		fmt.Printf("[ERR] Percona storage class could not be created\n")
		return
	}
	// display percona storage class
	display("percona storage class details", PerconaStorageClassDC)

	prompt("Create percona app")
	// create percona app
	create(PerconaAppCC)
	// display percona pvc
	display("percona pvc details", PerconaPVCDC)

	// verify if percona app is running
	isPARunning := verifyRunning(PerconaAppIRC)
	if !isPARunning {
		fmt.Printf("[ERR] Percona app is not running\n")
		return
	}

	// display openebs volume controller deployments
	display("openebs volume controller(s) details", OpenEBSVolumeCtrlDC)
	// display openebs volume replica deployments
	display("openebs volume replica(s) details", OpenEBSVolumeRepDC)

	// cleanup !!
	prompt("Proceed for cleanup")

	// delete percona app
	delete(PerconaAppDELC)

	// delete percona persistent volume claim
	delete(PerconaPVCDELC)

	// delete percona storage class
	delete(PerconaSCDELC)

	// delete host dir storage pool
	delete(HostDirSPDELC)

	// delete storage pool crd
	delete(StoragePoolCRDDELC)
}
