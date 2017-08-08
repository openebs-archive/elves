package operator

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/golang/glog"
	"github.com/openebs/elves/k8s-crd/sonar-operator/pkg/utils"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// Implements an Submarine's controller loop in a particular namespace.
// The controller makes use of an Informer resource to locally cache resources
// managed, and handle events on the resources.
type SubmarineController struct {
	// Baseline kubeconfig to use when communicating with the API.
	kubecfg *rest.Config

	// Clientset that has a REST client for each k8s API group.
	clientSet kubernetes.Interface

	// APIExtensions Clientset that has a REST client for each k8s API group.
	apiextensionsClientSet apiextensionsclient.Interface

	// REST client for the Submarine resource k8s API group (since its not an
	// official resource, there is no existing Clientset for it in k8s).
	restClient rest.Interface

	// Informer for all resources being watched by the operator.
	informer *SubmarineControllerInformer

	// The namespace where the operator is running.
	namespace string
}

// Implements an Informer for the resources being operated on: Submarines.
type SubmarineControllerInformer struct {
	// Store & controller for Submarine resources
	submarineStore      cache.Store
	submarineController cache.Controller
}

// Create a new Controller for the Submarine operator
func NewSubmarineController(kubeconfig, namespace) (
	*SubmarineController, error) {

	// Create the client config for use in creating the k8s API client
	// Use kubeconfig if given, otherwise use in-cluster
	kubecfg, err := utils.BuildKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// Create a new k8s API client from the kubeconfig
	clientSet, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, err
	}
	// Create a new k8s API client for API Extenstions from the kubeconfig
	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(kubecfg)
	if err != nil {
		return nil, err
	}

	// TODO - Use this block to check that CRD already exists
	// Create & register the Submarine resource as a CRD in the cluster, if it
	// doesn't exist
	//kind := reflect.TypeOf(Submarine{}).Name()
	//glog.V(2).Infof("Registering CRD: %s.%s | version: %s", CRDName, Domain, Version)
	//_, err = crd.CreateCustomResourceDefinition(
	//	apiextensionsClientSet,
	//	CRDName,
	//	Domain,
	//	kind,
	//	ResourceNamePlural,
	//	Version,
	//)
	//if err != nil {
	//	return nil, err
	//}

	// Discover or set the namespace in which this controller is running in
	if namespace == "" {
		if namespace = os.Getenv("POD_NAMESPACE"); namespace == "" {
			namespace = "default"
		}
	}

	// Create a new k8s REST API client for Submarines
	restClient, err := newSubmarineClient(kubecfg, namespace)
	if err != nil {
		return nil, err
	}

	// Create new SubmarineController
	sc := &SubmarineController{
		kubecfg:                kubecfg,
		clientSet:              clientSet,
		apiextensionsClientSet: apiextensionsClientSet,
		restClient:             restClient,
		namespace:              namespace,
	}

	// Create a new Informer for the SubmarineController
	sc.informer = sc.newSubmarineControllerInformer()

	return sc, nil
}

// Start the SubmarineController until stopped.
func (sc *SubmarineController) Start(stop <-chan struct{}) {
	// Don't let panics crash the process
	defer utilruntime.HandleCrash()

	glog.V(2).Infof("Starting Submarine controller...")
	sc.start(stop)

	// Block until stopped
	<-stop
}

// Start the controllers with the stop chan as required by Informers.
func (sc *SubmarineController) start(stop <-chan struct{}) {
	glog.V(2).Infof("Namespace: %s", sc.namespace)

	// Run controller for Submarine Informer and handle events via callbacks
	go sc.informer.submarineController.Run(stop)

	// Run the SubmarineController
	go sc.Run(stop)
}

// Informers are a combination of a local cache store to buffer the state of a
// given resource locally, and a controller to handle events through callbacks.
//
// Informers sync the APIServer's state of a resource with the local cache
// store.

// Creates a new Informer for the SubmarineController.
// An SubmarineController uses a set of Informers to watch and operate on
// Pods and Submarine resources in its control loop.
func (sc *SubmarineController) newSubmarineControllerInformer() *SubmarineControllerInformer {
	submarineStore, submarineController := sc.newSubmarineInformer()

	return &SubmarineControllerInformer{
		submarineStore:      submarineStore,
		submarineController: submarineController,
	}
}

// Create a new Informer on the Submarine resources in the cluster to
// track them.
func (sc *SubmarineController) newSubmarineInformer() (cache.Store, cache.Controller) {
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				// Retrieve an SubmarineList from the API
				result := &SubmarineList{}
				err := sc.restClient.Get().
					Resource(ResourceNamePlural).
					VersionedParams(&options, api.ParameterCodec).
					Do().
					Into(result)

				return result, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				// Watch the Submarines in the API
				return sc.restClient.Get().
					Prefix("watch").
					Namespace(sc.namespace).
					Resource(ResourceNamePlural).
					VersionedParams(&options, api.ParameterCodec).
					Watch()
			},
		},
		// The resource that the informer returns
		&Submarine{},
		// The sync interval of the informer
		5*time.Second,
		// Callback functions for add, delete & update events
		cache.ResourceEventHandlerFuncs{
			// AddFunc: func(o interface{}) {}
			UpdateFunc: sc.handleSubmarinesUpdate,
			// DeleteFunc: func(o interface{}) {}
		},
	)
}

// Callback for updates to an Submarine Informer
func (sc *SubmarineController) handleSubmarinesUpdate(oldObj, newObj interface{}) {
	// Make a copy of the object, to not mutate the original object from the local
	// cache store
	sub, err := CopyObjToSubmarine(newObj)
	if err != nil {
		glog.Errorf("Failed to copy Submarine object: %v", err)
		return
	}

	glog.V(2).Infof("Received update for Submarine: %s | "+
		"nation=%s ",
		sub.Metadata.Name,
		sub.Spec.Nation,
	)
}

// Run begins the SubmarineController.
func (sc *SubmarineController) Run(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			glog.V(2).Infof("Shutting down Submarine controller...")
			return
		default:
			sc.run()
			time.Sleep(2 * time.Second)
		}
	}
}

func (sc *SubmarineController) run() {
	glog.V(2).Infof("In SubmarineController loop...")

	// #########################################################################
	// Select the Submarine for the Namespace
	// #########################################################################

	// Use the first Submarine found as the operator currently
	// only supports a single Submarine per Namespace.
	// TODO (fix): This arbitrarily selects the first Submarine from the
	// list in the store. Order is not guaranteed.
	subListObj := sc.informer.submarineStore.List()
	subs, err := CopyObjToSubmarines(subListObj)
	if err != nil {
		glog.Errorf("Failed to copy object into Submarines: %v", err)
		return
	} else if len(subs) != 1 {
		glog.Errorf("No Submarines to list.")
		return
	}
	sub := subs[0]

	//TODO - Include the logic to raise an Event

}
