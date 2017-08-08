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

// Configure & create an k8s API REST client for the Submarine resource in the
// k8s cluster.
func newSubmarineClient(kubecfg *rest.Config, namespace string) (*rest.RESTClient, error) {
	// Update kubecfg to work with the Submarine's API group, using the kubecfg
	// param as a baseline.
	addSubmarineToKubeConfig(kubecfg, Domain, Version)

	// Add Submarine's API group to the k8s api.Scheme to provide it with the
	// capability of doing conversions or a deep-copy on an Submarine resource.
	addSubmarineToAPISchema(Domain, Version)

	// Create the k8s API REST client
	client, err := rest.RESTClientFor(kubecfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Configure the attributes for the kubecfg used in the Submarine API REST
// client.
func addSubmarineToKubeConfig(kubecfg *rest.Config, domain, version string) {
	groupversion := schema.GroupVersion{
		Group:   domain,
		Version: version,
	}

	// Set attributes in the kubecfg to reach and work with the
	// Submarine resource.
	kubecfg.GroupVersion = &groupversion
	kubecfg.APIPath = "/apis"
	kubecfg.ContentType = runtime.ContentTypeJSON
	kubecfg.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}
}

// Add the Submarine types to the api.Scheme for when needing to do type
// conversions or a deep-copy of an Submarine object.
func addSubmarineToAPISchema(domain, version string) {
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
			e.g. CopyObjToSubmarine()
	*/
	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			// AddKnownTypes registers all types passed in 'types' as being members of version 'version'.
			// All objects passed to types should be pointers to structs. The name that go reports for
			// the struct becomes the "kind" field when encoding. Version may not be empty - use the
			// APIVersionInternal constant if you have a type that does not have a formal version.
			scheme.AddKnownTypes(
				groupversion,
				&Submarine{},
				&SubmarineList{},
				&metav1.ListOptions{},
				&metav1.DeleteOptions{},
			)
			return nil
		})

	// AddToScheme applies all the stored functions to the scheme.A non-nil error
	// indicates that one function failed and the attempt was abandoned.
	schemeBuilder.AddToScheme(api.Scheme)
}

// Create a deep-copy of an Submarine object
func CopyObjToSubmarine(obj interface{}) (*Submarine, error) {
	objCopy, err := api.Scheme.Copy(obj.(*Submarine))
	if err != nil {
		return nil, err
	}

	am := objCopy.(*Submarine)
	if am.Metadata.Annotations == nil {
		am.Metadata.Annotations = make(map[string]string)
	}
	return am, nil
}

// Attempt to deep copy an empty interface into an SubmarineList.
func CopyObjToSubmarines(obj []interface{}) ([]Submarine, error) {
	subs := []Submarine{}

	for _, o := range obj {
		sub, err := CopyObjToSubmarine(o)
		if err != nil {
			glog.Errorf("Failed to copy submarine object for subsList: %v", err)
			return nil, err
		}
		subs = append(subs, *sub)
	}

	return subs, nil
}

// #############################################################################

/*
 Note: The following code is boilerplate code needed to satisfy the
 Submarine as a resource in the cluster in terms of how it expects CRD's to
 be created, operate and used.
*/

// #############################################################################

// Required to satisfy Object interface
func (sub *Submarine) GetObjectKind() schema.ObjectKind {
	return &sub.TypeMeta
}

// Required to satisfy ObjectMetaAccessor interface
func (sub *Submarine) GetObjectMeta() metav1.Object {
	return &sub.Metadata
}

// Required to satisfy Object interface
func (subs *SubmarineList) GetObjectKind() schema.ObjectKind {
	return &subs.TypeMeta
}

// Required to satisfy ListMetaAccessor interface
func (subs *SubmarineList) GetListMeta() metav1.List {
	return &subs.Metadata
}

// #############################################################################

/*
 Note: The following code is used only to work around a known problem
 with third-party resources and ugorji. If/when these issues are resolved,
 the code below should no longer be required.
*/

// #############################################################################

type SubmarineListCopy SubmarineList
type SubmarineCopy Submarine

func (sub *Submarine) UnmarshalJSON(data []byte) error {
	tmp := SubmarineCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := Submarine(tmp)
	*sub = tmp2
	return nil
}

func (subs *SubmarineList) UnmarshalJSON(data []byte) error {
	tmp := SubmarineListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmp2 := SubmarineList(tmp)
	*subs = tmp2
	return nil
}

// #############################################################################

/*
 Note: The rest of the following code is not currently used in this project.
 It is provided to serve as additional references & examples on
 communicating with the cluster API.
*/

// #############################################################################

func NewSubmarine(name string, nation string) *Submarine {
	return &Submarine{
		Metadata: metav1.ObjectMeta{
			Name: name,
		},
		Spec: SubmarineSpec{
			Nation: nation,
		},
	}
}

// Instantiate an Submarine in the cluster.
func (sub *Submarine) Instantiate(kubecfg, namespace string) error {
	// Create the client config. Use the kubecfg if provided, else
	// assume we're using automatic in-cluster auth.
	k, err := utils.BuildKubeConfig(kubecfg)
	if err != nil {
		return err
	}

	sub.testCreate(k, namespace)
	return nil
}

// Test the creation of an Submarine in the cluster
func (sub *Submarine) testCreate(kubecfg *rest.Config, namespace string) {
	// Create a new k8s API client for Submarines
	client, err := newSubmarineClient(kubecfg, namespace)
	if err != nil {
		panic(err)
	}

	err = client.Get().
		Resource(ResourceNamePlural).
		Namespace(namespace).
		Name(am.Metadata.Name).
		Do().Into(am)

	if err != nil {
		// Create the Submarine if it doesn't already exist
		if errors.IsNotFound(err) {
			var result Submarine
			// Attempt to POST the Submarine to the APIServer every 1 second, exiting
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
		// Else, it already exists, GET / retreive the Submarine
	} else {
		fmt.Printf("GET: %#v\n", sub)
	}
}

// List all Submarines in the cluster.
func ListSubmarinesWithClient(kubecfg *rest.Config, namespace string) {
	// Create a new k8s API client for Submarines
	client, err := newSubmarineClient(kubecfg, namespace)
	if err != nil {
		panic(err)
	}

	// Fetch a list of our CRDs
	submarineList := SubmarineList{}
	err = client.Get().Resource(ResourceNamePlural).Do().Into(&submarineList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", submarineList)
}
