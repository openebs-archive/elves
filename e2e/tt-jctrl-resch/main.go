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
	"fmt"
	"net/url"
	"time"
  "encoding/json"

	mayav1 "github.com/openebs/maya/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	k8sclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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

// defaultVolSpecs returns the default specifications of openebs jiva based volume
func defaultVolSpecs() string {
	return `{"kind":"PersistentVolumeClaim","apiVersion":"v1","metadata":{"name":"my-jiva-vsm"}}`
}

// getVolSvcHttpStatusByName returns the http status code of openebs volume's
// controller service
func getVolSvcHTTPStatusByName(volServiceName string, k8sclient k8sclientset.Interface) (int, error) {
	volSvcIP, err := getK8sSvcIPByName(volServiceName, k8sclient)
	if err != nil {
		return 500, err
	}
	ovrc, err := volRESTClient(volSvcIP)
	if err != nil {
		return 500, err
	}

	// http get
	result := ovrc.Get().Do()
	var statusCode int
	result.StatusCode(&statusCode)
	return statusCode, nil
}

func createVolByMApiService(mAPIServiceName string, defaultVolSpecs string, k8sclient k8sclientset.Interface) (string, error) {
	// Get the maya api service IP address
	mAPISvcIP, err := getK8sSvcIPByName(mAPIServiceName, k8sclient)
	if err != nil {
		return "", err
	}
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
		return "", fmt.Errorf("Missing name in openebs volume")
	}
	return pv.Name, nil
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	k8sclient, err := k8sclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// business logic here
	// 1. create a openebs volume consisting of jiva controller & replica(s)
	volName, err := createVolByMApiService("maya-apiserver-service", defaultVolSpecs(), k8sclient)
	if err != nil {
	  panic(err.Error())
	}
	// 2. initiate a endless loop that captures the jiva volume's service status
	for {
	  status, err := getVolSvcHTTPStatusByName(volName + "-ctrl-svc", k8sclient)
	  if err != nil {
	    panic(err.Error())
	  }
	  fmt.Printf("%s: %d\n", time.Now().String(), status)
		time.Sleep(1 * time.Second)
	}
	// 3. outside of the program: get the node with jiva controller & shut it down
	// 4. outside of the program: check the pod logs & estimate the time taken for rescheduling
}
