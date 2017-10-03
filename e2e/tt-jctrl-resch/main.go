/*
Copyright 2016 The Kubernetes Authors.
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
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	mayav1 "github.com/openebs/maya/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"

	k8sclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// DefaultSchemeGroupVersion is group version used to register objects
var DefaultSchemeGroupVersion = schema.GroupVersion{Group: "", Version: "v1"}

// mAPIReqContentConfig
func mAPIReqContentConfig() rest.ContentConfig {
	gvCopy := DefaultSchemeGroupVersion
	return rest.ContentConfig{
		ContentType:          "application/json",
		GroupVersion:         &gvCopy,
		NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: scheme.Codecs},
	}
}

// mAPIRESTClient returns the REST client capable of communicating with
// maya api service
func mAPIRESTClient(mAPISvcIP string) (*rest.RESTClient, error) {
	baseURL, err := url.Parse("http://" + mAPISvcIP + ":5656")
	if err != nil {
		return nil, err
	}

	client, err := rest.NewRESTClient(baseURL, "/latest/", mAPIReqContentConfig(), 0, 0, nil, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// volRESTClient returns the REST client capable of communicating with
// openebs volume's jiva controller service
func volRESTClient(volSvcIP string) (*rest.RESTClient, error) {
	baseURL, err := url.Parse("http://" + volSvcIP + ":9501")
	if err != nil {
		return nil, err
	}

	client, err := rest.NewRESTClient(baseURL, "", mAPIReqContentConfig(), 0, 0, nil, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// getK8sSvcIPByName fetches the Cluster IP Address of any Kubernetes service.
//
// TODO
//  k8sclient should be a property from a struct e.g. similar to framework of k8s e2e
// This future struct should be passed around. This struct should be populated
// based on kubeconfig etc.
func getK8sSvcIPByName(svcName string, k8sclient k8sclientset.Interface) (string, error) {
	svcClient := k8sclient.CoreV1().Services("default")

	svc, err := svcClient.Get(svcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if svc == nil {
		return "", fmt.Errorf("service '%s' not found", svcName)
	}

	return svc.Spec.ClusterIP, nil
}

//
func getK8sDeploymentByName(deployName string, k8sclient k8sclientset.Interface) (*v1beta1.Deployment, error) {
	deployClient := k8sclient.ExtensionsV1beta1().Deployments("default")

	deploy, err := deployClient.Get(deployName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if deploy == nil {
		return nil, fmt.Errorf("deploy '%s' not found", deployName)
	}

	return deploy, nil
}

// getControllerPods fetches the Controller Pods
func getOVControllerPods(volumeName string, k8sclient k8sclientset.Interface) (*v1.PodList, error) {

	podClient := k8sclient.CoreV1().Pods("default")

	// filter the VSM Controller Pod(s)
	pOpts := metav1.ListOptions{
		// A list of comma separated key=value filters will filter the
		// OpenEBS Volume Controller Pod(s)
		LabelSelector: string(mayav1.VSMSelectorKeyEquals) + volumeName + "," + string(mayav1.ControllerSelectorKeyEquals) + string(mayav1.JivaControllerSelectorValue),
	}

	cps, err := podClient.List(pOpts)
	if err != nil {
		return nil, err
	}

	return cps, nil
}

// WaitForDeploymentReady will wait for the deployment to be ready
func WaitForDeploymentReady(timeout time.Duration, deployName string, k8sclient k8sclientset.Interface) error {

	// POLL until the deployment is ready
	// Check this every X seconds till the timeout
	return wait.Poll(5*time.Second, timeout,
		func() (bool, error) {
			deploy, err := getK8sDeploymentByName(deployName, k8sclient)
			if err != nil {
				return false, err
			}
			if deploy.Status.UnavailableReplicas != 0 {
				return false, fmt.Errorf("Deployment '%s' has '%d' unavailable replica(s)", deployName, deploy.Status.UnavailableReplicas)
			}

			return true, nil
		})
}

// reschOVControllerByName will initiate a reschedule of openebs volume controller
// deployment by deleting its Pod(s) one-by-one.
func reschOVControllerByName(volumeName string, k8sclient k8sclientset.Interface) error {

	cps, err := getOVControllerPods(volumeName, k8sclient)
	if err != nil {
		return err
	}

	podClient := k8sclient.CoreV1().Pods("default")

	for _, cp := range cps.Items {
		// Check if pod is running else error out
		if cp.Status.Phase != v1.PodRunning {
			return fmt.Errorf("Pod '%s' is not running", cp.Name)
		}

		// Delete the K8s Pod
		// This will trigger a reschedule as this Pod is managed by a K8s Deployment
		err := podClient.Delete(cp.Name, &metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// defaultVolSpecs returns the default specifications of openebs jiva based volume
func defaultVolSpecs() string {
	return `{"kind":"PersistentVolumeClaim","apiVersion":"v1","metadata":{"name":"e2e-tt-jctrl-resch-vol"}}`
}

// getVolSvcHttpStatusByName returns the http status code of openebs volume's
// controller service
func getVolSvcHTTPStatusByName(volServiceName string, k8sclient k8sclientset.Interface) (string, int, error) {
	volSvcIP, err := getK8sSvcIPByName(volServiceName, k8sclient)
	if err != nil {
		return "", 500, err
	}
	ovrc, err := volRESTClient(volSvcIP)
	if err != nil {
		return "", 500, err
	}

	// http get
	result := ovrc.Get().Do()
	var statusCode int
	result.StatusCode(&statusCode)

	return time.Now().String(), statusCode, nil
}

//
func repeatGetVolSvcHTTPStatusByName(repeat int, intervalSecs int, volServiceName string, k8sclient k8sclientset.Interface) error {

	if repeat <= 0 {
		return fmt.Errorf("Invalid repeat count '%d'. It should be a positive number", repeat)
	}

	if intervalSecs <= 0 {
		return fmt.Errorf("Invalid interval '%d'. It should be a positive number", intervalSecs)
	}

	fmt.Printf("Start capturing http status of volume service\n")

	// repeat in an interval of X seconds
	const delay = 1 * 1000 * time.Millisecond
	for i := 1; i <= repeat; i++ {
		tstamp, scode, err := getVolSvcHTTPStatusByName(volServiceName, k8sclient)
		if err != nil {
			return fmt.Errorf("Counter '#%d' resulted in error '%s'", i, err.Error())
		}

		// TODO
		// export these as events to Prometheus
		// Probably decorate this method with a method that exclusively exports
		// to Prometheus
		fmt.Printf("Counter: '#%d', TimeStamp: '%s', HttpStatus: '%d'\n", i, tstamp, scode)

		// Sleep a bit
		time.Sleep(delay)
	}

	fmt.Printf("Completed capturing http status of volume service\n")

	return nil
}

func createVolByMApiService(mAPIServiceName string, defaultVolSpecs string, k8sclient k8sclientset.Interface) (string, error) {
	// Get the maya api service IP address
	mAPISvcIP, err := getK8sSvcIPByName(mAPIServiceName, k8sclient)
	if err != nil {
		return "", err
	}
	// Build the REST client to maya api service
	mapiRC, err := mAPIRESTClient(mAPISvcIP)
	if err != nil {
		return "", err
	}

	// http post to create
	contents, err := mapiRC.Post().AbsPath("/latest/volumes/").Body([]byte(defaultVolSpecs)).SetHeader("Accept", "application/json").DoRaw()
	if err != nil {
		return "", err
	}
	// set against the proper type !!
	var pv mayav1.PersistentVolume
	err = json.Unmarshal(contents, &pv)
	if pv.Name == "" {
		return "", fmt.Errorf("Name could not be determined in the created openebs volume")
	}
	return pv.Name, nil
}

//
func deleteVolByMApiService(mAPIServiceName string, volName string, k8sclient k8sclientset.Interface) error {
	// Get the maya api service IP address
	mAPISvcIP, err := getK8sSvcIPByName(mAPIServiceName, k8sclient)
	if err != nil {
		return err
	}
	// Build the REST client to maya api service
	mapiRC, err := mAPIRESTClient(mAPISvcIP)
	if err != nil {
		return err
	}

	// http GET to delete ?? m-apiserver this is BAD !!!
	_, err = mapiRC.Get().AbsPath("/latest/volumes/delete/" + volName).DoRaw()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Signal for goroutine
	// Used when goroutine logic encounters an error
	errGetStatusSig := make(chan error)

	// creates the in-cluster config
	// rest.InClusterConfig() uses the Service Account token mounted inside the
	// Pod at /var/run/secrets/kubernetes.io/serviceaccount path.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	k8sclient, err := k8sclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Business Logic Starts Here
	// 1. create a openebs volume consisting of jiva controller & replica(s)
	volName, err := createVolByMApiService("maya-apiserver-service", defaultVolSpecs(), k8sclient)
	if err != nil {
		panic(err.Error())
	}

	// 2.1 Fetch the openebs volume controller deployment
	deploy, err := getK8sDeploymentByName(volName+"-ctrl", k8sclient)
	if err != nil {
		panic(err.Error())
	}
	// 2.2 Wait till openebs volume controller deployment is Ready
	err = WaitForDeploymentReady(3*time.Minute, deploy.Name, k8sclient)
	if err != nil {
		panic(err.Error())
	}

	// 3. Repeatedly capture the openebs volume controller service status for X times
	// in a goroutine. The (X times * (Y seconds interval + execution time))
	// should be good enough to capture the status both before & after the controller
	// reschedule. The change of status from 200 to 0 & back to 200 will give us
	// the time taken for controller reschedule.
	go func() {
		err = repeatGetVolSvcHTTPStatusByName(50, 1, volName+"-ctrl-svc", k8sclient)
		errGetStatusSig <- err
	}()

	// 4.1 Wait for couple of seconds before triggering the re-schedule
	time.Sleep(5 * 1000 * time.Millisecond)

	// 4.2 reschedule openebs volume controller
	err = reschOVControllerByName(volName, k8sclient)
	if err != nil {
		panic(err.Error())
	}

	// TODO
	// Move away from manual checking of log files -> to -> exporting events to
	// Prometheus. Automation is the key. Use Prometheus as reporting tool.
	// Later use Prometheus as eventing mechansim.

	// 4.3 Do this manually till above automation is done
	// outside of the program: check the pod logs & estimate the time taken
	// for rescheduling

	// 5.1 Wait till the finish of goroutine w.r.t repeatGetVolSvcHTTPStatusByName
	err = <-errGetStatusSig
	if err != nil {
		panic(err.Error())
	}

	// 5.2 Delete the openebs volume as we are done with the testing
	err = deleteVolByMApiService("maya-apiserver-service", volName, k8sclient)
	if err != nil {
		panic(err.Error())
	}
}
