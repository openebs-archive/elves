package operator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/metral/memhog-operator/pkg/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	// To authenticate against GKE clusters
	//_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// #############################################################################

/*
 Note: The k8s.io/client-go lib's REST client.Get() and client.List()
 provide a means of accessing & working with a particular resource in the
 cluster.
 However, using client.Get() and client.List() on the REST client can
 become expensive if used multiple times.

 A user can benefit working with the cluster through its client by using
 a local cache store and event watches for better performance.
 This optimization is suggested and can be implemented with an:
 	- Informer: A local cache store & controller for state event handling on a
 	resource, that sync’s with the APIServer's state.
  See https://github.com/kubernetes/client-go/blob/v2.0.0/tools/cache/controller.go#L201-L221
 	- SharedInformer: A single, optimized local cache store & controller for
 	state event handling on multiple resources, syncing all stores &
 	controllers with the APIServer's state.
  See https://github.com/kubernetes/client-go/blob/v2.0.0/tools/cache/shared_informer.go#L31-L39
*/

// #############################################################################

// Configure & create an k8s API REST client for the AppMonitor resource in the
// k8s cluster.
func newAppMonitorClient(kubecfg *rest.Config, namespace string) (*rest.RESTClient, error) {
	// Update kubecfg to work with the AppMonitor's API group, using the kubecfg
	// param as a baseline.
	addAppMonitorToKubeConfig(kubecfg, Domain, Version)

	// Add AppMonitor's API group to the k8s api.Scheme to provide it with the
	// capability of doing conversions or a deep-copy on an AppMonitor resource.
	addAppMonitorToAPISchema(Domain, Version)

	// Create the k8s API REST client
	client, err := rest.RESTClientFor(kubecfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Configure the attributes for the kubecfg used in the AppMonitor API REST
// client.
func addAppMonitorToKubeConfig(kubecfg *rest.Config, domain, version string) {
	groupversion := schema.GroupVersion{
		Group:   domain,
		Version: version,
	}

	// Set attributes in the kubecfg to reach and work with the
	// AppMonitor resource.
	kubecfg.GroupVersion = &groupversion
	kubecfg.APIPath = "/apis"
	kubecfg.ContentType = runtime.ContentTypeJSON
	kubecfg.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}
}

// Add the AppMonitor types to the api.Scheme for when needing to do type
// conversions or a deep-copy of an AppMonitor object.
func addAppMonitorToAPISchema(domain, version string) {
	groupversion := schema.GroupVersion{
		Group:   domain,
		Version: version,
	}

	/*
		 Scheme defines methods for serializing and deserializing API objects, a type
		 registry for converting group, version, and kind information to and from Go
		 schemas, and mappings between Go schemas of different versions. A scheme is the
		 foundation for a versioned API and versioned configuration over time.

		 In a Scheme, a Type is a particular Go struct, a Version is a point-in-time
		 identifier for a particular representation of that Type (typically backwards
		 compatible), a Kind is the unique name for that Type within the Version, and a
		 Group identifies a set of Versions, Kinds, and Types that evolve over time. An
		 Unversioned Type is one that is not yet formally bound to a type and is promised
		 to be backwards compatible (effectively a "v1" of a Type that does not expect
		 to break in the future).

		 SchemeBuilder collects functions that add things to a scheme. It's to
		 allow code to compile without explicitly referencing generated types.
		 You should declare one in each package that will have generated deep-copy
		 or conversion functions.

		 Create a schemeBuilder that will ultimately add / register the
		 following types for groupversion into the api.Scheme used when performing
		 a deep-copy of an opject as to not mutate the original object.
			e.g. CopyObjToAppMonitor()
	*/
	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			// AddKnownTypes registers all types passed in 'types' as being members of version 'version'.
			// All objects passed to types should be pointers to structs. The name that go reports for
			// the struct becomes the "kind" field when encoding. Version may not be empty - use the
			// APIVersionInternal constant if you have a type that does not have a formal version.
			scheme.AddKnownTypes(
				groupversion,
				&AppMonitor{},
				&AppMonitorList{},
				&metav1.ListOptions{},
				&metav1.DeleteOptions{},
			)
			return nil
		})

	// AddToScheme applies all the stored functions to the scheme.A non-nil error
	// indicates that one function failed and the attempt was abandoned.
	schemeBuilder.AddToScheme(api.Scheme)
}

// Create a deep-copy of an AppMonitor object
func CopyObjToAppMonitor(obj interface{}) (*AppMonitor, error) {
	objCopy, err := api.Scheme.Copy(obj.(*AppMonitor))
	if err != nil {
		return nil, err
	}

	am := objCopy.(*AppMonitor)
	if am.Metadata.Annotations == nil {
		am.Metadata.Annotations = make(map[string]string)
	}
	return am, nil
}

// Attempt to deep copy an empty interface into an AppMonitorList.
func CopyObjToAppMonitors(obj []interface{}) ([]AppMonitor, error) {
	ams := []AppMonitor{}

	for _, o := range obj {
		am, err := CopyObjToAppMonitor(o)
		if err != nil {
			glog.Errorf("Failed to copy pod object for podList: %v", err)
			return nil, err
		}
		ams = append(ams, *am)
	}

	return ams, nil
}

// #############################################################################

/*
 Note: The following code is boilerplate code needed to satisfy the
 AppMonitor as a resource in the cluster in terms of how it expects CRD's to
 be created, operate and used.
*/

// #############################################################################

// Required to satisfy Object interface
func (am *AppMonitor) GetObjectKind() schema.ObjectKind {
	return &am.TypeMeta
}

// Required to satisfy ObjectMetaAccessor interface
func (am *AppMonitor) GetObjectMeta() metav1.Object {
	return &am.Metadata
}

// Required to satisfy Object interface
func (aml *AppMonitorList) GetObjectKind() schema.ObjectKind {
	return &aml.TypeMeta
}

// Required to satisfy ListMetaAccessor interface
func (aml *AppMonitorList) GetListMeta() metav1.List {
	return &aml.Metadata
}

// #############################################################################

/*
 Note: The following code is used only to work around a known problem
 with third-party resources and ugorji. If/when these issues are resolved,
 the code below should no longer be required.
*/

// #############################################################################

type AppMonitorListCopy AppMonitorList
type AppMonitorCopy AppMonitor

func (am *AppMonitor) UnmarshalJSON(data []byte) error {
	tmp := AppMonitorCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := AppMonitor(tmp)
	*am = tmp2
	return nil
}

func (aml *AppMonitorList) UnmarshalJSON(data []byte) error {
	tmp := AppMonitorListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := AppMonitorList(tmp)
	*aml = tmp2
	return nil
}

// #############################################################################

/*
 Note: The rest of the following code is not currently used in this project.
 It is provided to serve as additional references & examples on
 communicating with the cluster API.
*/

// #############################################################################

func NewAppMonitor(name string, memThresholdPercent, memMultiplier float64) *AppMonitor {
	return &AppMonitor{
		Metadata: metav1.ObjectMeta{
			Name: name,
		},
		Spec: AppMonitorSpec{
			MemThresholdPercent: memThresholdPercent,
			MemMultiplier:       memMultiplier,
		},
	}
}

// Instantiate an AppMonitor in the cluster.
func (am *AppMonitor) Instantiate(kubecfg, namespace string) error {
	// Create the client config. Use the kubecfg if provided, else
	// assume we're using automatic in-cluster auth.
	k, err := utils.BuildKubeConfig(kubecfg)
	if err != nil {
		return err
	}

	am.testCreate(k, namespace)
	return nil
}

// Test the creation of an AppMonitor in the cluster
func (am *AppMonitor) testCreate(kubecfg *rest.Config, namespace string) {
	// Create a new k8s API client for AppMonitors
	client, err := newAppMonitorClient(kubecfg, namespace)
	if err != nil {
		panic(err)
	}

	err = client.Get().
		Resource(ResourceNamePlural).
		Namespace(namespace).
		Name(am.Metadata.Name).
		Do().Into(am)

	if err != nil {
		// Create the AppMonitor if it doesn't already exist
		if errors.IsNotFound(err) {
			var result AppMonitor
			// Attempt to POST the AppMonitor to the APIServer every 1 second, exiting
			// when successfuly created.
			for err != nil {
				err = client.Post().
					Resource(ResourceNamePlural).
					Namespace(namespace).
					Body(am).
					Do().Into(&result)

				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				fmt.Printf("CREATED: %#v\n", result)
			}
		} else {
			panic(err)
		}
		// Else, it already exists, GET / retreive the AppMonitor
	} else {
		fmt.Printf("GET: %#v\n", am)
	}
}

// List all AppMonitors in the cluster.
func ListAppMonitorsWithClient(kubecfg *rest.Config, namespace string) {
	// Create a new k8s API client for AppMonitors
	client, err := newAppMonitorClient(kubecfg, namespace)
	if err != nil {
		panic(err)
	}

	// Fetch a list of our CRDs
	appmonitorList := AppMonitorList{}
	err = client.Get().Resource(ResourceNamePlural).Do().Into(&appmonitorList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", appmonitorList)
}
