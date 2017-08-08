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
	"github.com/metral/memhog-operator/pkg/operator/crd"
	"github.com/metral/memhog-operator/pkg/utils"
	prometheusClient "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// Implements an AppMonitor's controller loop in a particular namespace.
// The controller makes use of an Informer resource to locally cache resources
// managed, and handle events on the resources.
type AppMonitorController struct {
	// Baseline kubeconfig to use when communicating with the API.
	kubecfg *rest.Config

	// Clientset that has a REST client for each k8s API group.
	clientSet kubernetes.Interface

	// APIExtensions Clientset that has a REST client for each k8s API group.
	apiextensionsClientSet apiextensionsclient.Interface

	// REST client for the AppMonitor resource k8s API group (since its not an
	// official resource, there is no existing Clientset for it in k8s).
	restClient rest.Interface

	// Informer for all resources being watched by the operator.
	informer *AppMonitorControllerInformer

	// The namespace where the operator is running.
	namespace string

	// The address of the Prometheus service
	// e.g. "http://prometheus.tectonic-system:9090"
	prometheusAddr string
}

// Implements an Informer for the resources being operated on: Pods &
// AppMonitors.
type AppMonitorControllerInformer struct {
	// Store & controller for Pod resources
	podStore      cache.Store
	podController cache.Controller

	// Store & controller for AppMonitor resources
	appMonitorStore      cache.Store
	appMonitorController cache.Controller
}

// Create a new Controller for the AppMonitor operator
func NewAppMonitorController(kubeconfig, namespace, prometheusAddr string) (
	*AppMonitorController, error) {

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

	// Create & register the AppMonitor resource as a CRD in the cluster, if it
	// doesn't exist
	kind := reflect.TypeOf(AppMonitor{}).Name()
	glog.V(2).Infof("Registering CRD: %s.%s | version: %s", CRDName, Domain, Version)
	_, err = crd.CreateCustomResourceDefinition(
		apiextensionsClientSet,
		CRDName,
		Domain,
		kind,
		ResourceNamePlural,
		Version,
	)
	if err != nil {
		return nil, err
	}

	// Discover or set the namespace in which this controller is running in
	if namespace == "" {
		if namespace = os.Getenv("POD_NAMESPACE"); namespace == "" {
			namespace = "default"
		}
	}

	// Create a new k8s REST API client for AppMonitors
	restClient, err := newAppMonitorClient(kubecfg, namespace)
	if err != nil {
		return nil, err
	}

	// Create new AppMonitorController
	amc := &AppMonitorController{
		kubecfg:                kubecfg,
		clientSet:              clientSet,
		apiextensionsClientSet: apiextensionsClientSet,
		restClient:             restClient,
		namespace:              namespace,
		prometheusAddr:         prometheusAddr,
	}

	// Create a new Informer for the AppMonitorController
	amc.informer = amc.newAppMonitorControllerInformer()

	return amc, nil
}

// Start the AppMonitorController until stopped.
func (amc *AppMonitorController) Start(stop <-chan struct{}) {
	// Don't let panics crash the process
	defer utilruntime.HandleCrash()

	glog.V(2).Infof("Starting AppMonitor controller...")
	amc.start(stop)

	// Block until stopped
	<-stop
}

// Start the controllers with the stop chan as required by Informers.
func (amc *AppMonitorController) start(stop <-chan struct{}) {
	glog.V(2).Infof("Namespace: %s", amc.namespace)

	// Run controller for Pod Informer and handle events via callbacks
	go amc.informer.podController.Run(stop)

	// Run controller for AppMonitor Informer and handle events via callbacks
	go amc.informer.appMonitorController.Run(stop)

	// Run the AppMonitorController
	go amc.Run(stop)
}

// Informers are a combination of a local cache store to buffer the state of a
// given resource locally, and a controller to handle events through callbacks.
//
// Informers sync the APIServer's state of a resource with the local cache
// store.

// Creates a new Informer for the AppMonitorController.
// An AppMonitorController uses a set of Informers to watch and operate on
// Pods and AppMonitor resources in its control loop.
func (amc *AppMonitorController) newAppMonitorControllerInformer() *AppMonitorControllerInformer {
	podStore, podController := amc.newPodInformer()
	appMonitorStore, appMonitorController := amc.newAppMonitorInformer()

	return &AppMonitorControllerInformer{
		podStore:             podStore,
		podController:        podController,
		appMonitorStore:      appMonitorStore,
		appMonitorController: appMonitorController,
	}
}

// Create a new Informer on the Pod resources in the cluster to track them.
func (amc *AppMonitorController) newPodInformer() (cache.Store, cache.Controller) {
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return amc.clientSet.CoreV1().Pods(amc.namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return amc.clientSet.CoreV1().Pods(amc.namespace).Watch(options)
			},
		},
		// The resource that the informer returns
		&v1.Pod{},
		// The sync interval of the informer
		5*time.Second,
		// Callback functions for add, delete & update events
		cache.ResourceEventHandlerFuncs{
			// AddFunc: func(o interface{}) {}
			UpdateFunc: amc.handlePodsUpdate,
			// DeleteFunc: func(o interface{}) {}
		},
	)
}

//TODO: Run controller loop
// - Watch for Pods that contain the 'app-monitor.kubedemo.com: true'
// - Watch for any new AppMonitor resources in the current namespace.
// annotation to indicate that an AppMonitor should operate on it.
// - Poll annotated Pod's heap size rate for an interval of time from Prometheus
// - Compare Pod heap size to Pod memory limit using the memThresholdPercent

// Callback for updates to a Pod Informer
func (amc *AppMonitorController) handlePodsUpdate(oldObj, newObj interface{}) {
	// Make a copy of the object, to not mutate the original object from the local
	// cache store
	pod, err := utils.CopyObjToPod(newObj)
	if err != nil {
		glog.Errorf("Failed to copy Pod object: %v", err)
		return
	}

	if _, ok := pod.Annotations[AppMonitorAnnotation]; !ok {
		return
	}

	glog.V(2).Infof("Received update for annotated Pod: %s | Annotations: %s", pod.Name, pod.Annotations)
}

// Create a new Informer on the AppMonitor resources in the cluster to
// track them.
func (amc *AppMonitorController) newAppMonitorInformer() (cache.Store, cache.Controller) {
	return cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				// Retrieve an AppMonitorList from the API
				result := &AppMonitorList{}
				err := amc.restClient.Get().
					Resource(ResourceNamePlural).
					VersionedParams(&options, api.ParameterCodec).
					Do().
					Into(result)

				return result, err
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				// Watch the AppMonitors in the API
				return amc.restClient.Get().
					Prefix("watch").
					Namespace(amc.namespace).
					Resource(ResourceNamePlural).
					VersionedParams(&options, api.ParameterCodec).
					Watch()
			},
		},
		// The resource that the informer returns
		&AppMonitor{},
		// The sync interval of the informer
		5*time.Second,
		// Callback functions for add, delete & update events
		cache.ResourceEventHandlerFuncs{
			// AddFunc: func(o interface{}) {}
			UpdateFunc: amc.handleAppMonitorsUpdate,
			// DeleteFunc: func(o interface{}) {}
		},
	)
}

// Callback for updates to an AppMonitor Informer
func (amc *AppMonitorController) handleAppMonitorsUpdate(oldObj, newObj interface{}) {
	// Make a copy of the object, to not mutate the original object from the local
	// cache store
	am, err := CopyObjToAppMonitor(newObj)
	if err != nil {
		glog.Errorf("Failed to copy AppMonitor object: %v", err)
		return
	}

	glog.V(2).Infof("Received update for AppMonitor: %s | "+
		"memThresholdPercent=%.2f | memMultiplier=%.2f",
		am.Metadata.Name,
		am.Spec.MemThresholdPercent,
		am.Spec.MemMultiplier,
	)
}

// Run begins the AppMonitorController.
func (amc *AppMonitorController) Run(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			glog.V(2).Infof("Shutting down AppMonitor controller...")
			return
		default:
			amc.run()
			time.Sleep(2 * time.Second)
		}
	}
}

func (amc *AppMonitorController) run() {
	glog.V(2).Infof("In AppMonitorController loop...")

	// #########################################################################
	// Select the Pods annotated to be managed by an AppMonitor
	// #########################################################################

	// List the pods as a []interface{} from the local cache store, representing
	// a PodList.
	podListObj := amc.informer.podStore.List()
	var err error

	// Make a copy of the object, to not mutate the original object from the local
	// cache store.
	pods, err := utils.CopyObjToPods(podListObj)
	if err != nil {
		glog.Errorf("Failed to copy object into Pods: %v", err)
		return
	}

	// Select Pods that are annotated to be monitored by an AppMonitor.
	var annotatedPods []v1.Pod
	annotatedPods, err = utils.SelectAnnotatedPods(pods, AppMonitorAnnotation)
	if err != nil {
		glog.Errorf("Failed to select annotated pods: %v", err)
		return
	}

	// #########################################################################
	// Select the AppMonitor for the Namespace
	// #########################################################################

	// Use the first AppMonitor found as the operator currently
	// only supports a single AppMonitor per Namespace.
	// TODO (fix): This arbitrarily selects the first AppMonitor from the
	// list in the store. Order is not guaranteed.
	amListObj := amc.informer.appMonitorStore.List()
	ams, err := CopyObjToAppMonitors(amListObj)
	if err != nil {
		glog.Errorf("Failed to copy object into AppMonitors: %v", err)
		return
	} else if len(ams) != 1 {
		glog.Errorf("No AppMonitors to list.")
		return
	}
	am := ams[0]

	// #########################################################################
	// Iterate on the annotated Pods, searching for Pods in need of vertical
	// scaling / redeployment.
	// #########################################################################

	for _, pod := range annotatedPods {
		glog.V(2).Infof("Iterating on Annotated Pod from Store PodList: %s", pod.Name)
		glog.V(2).Infof("Iterating on AppMonitor from Store AppMonitorList: %s",
			am.Metadata.Name)

		// Pull metrics off Pod & AppMonitor to determine if Pod needs a redeploy.
		container := pod.Spec.Containers[0]
		podLimits, _ := container.Resources.Limits["memory"]
		podLimitsBytes, _ := podLimits.AsInt64()

		// Query Prometheus for the current bytes for the Pod

		// Create a new Prometheus Client
		c, err := prometheusClient.NewClient(
			prometheusClient.Config{
				Address: amc.prometheusAddr,
			},
		)
		if err != nil {
			glog.Errorf("Failed to create client to Prometheus API: %v", err)
			return
		}

		promClient := prometheus.NewAPI(c)
		queryString := `container_memory_usage_bytes
			{
				namespace="default",
			  pod_name="%s",
				container_name="memhog"
			}
			`
		query := fmt.Sprintf(queryString, pod.Name)
		glog.V(2).Infof("Prometheus Query: %s", query)

		rawResults, err := queryPrometheus(promClient, query)
		if err != nil {
			glog.Errorf("Failed to query the Prometheus API: %v", err)
			return
		}
		glog.V(2).Infof("Prometheus query raw results: %s", rawResults)

		results := getMatrixValuesFromResults(rawResults)
		if len(results) == 0 {
			glog.V(2).Infof("Prometheus query results: empty")
			return
		}

		// Arbitrarily choose the first metric values returned in a
		// possible series of metrics
		r := results[0]

		// Retrieve the metric in the (time, metric) tuple
		for len(r.Values) == 0 {
			time.Sleep(500 * time.Millisecond)
		}
		val := r.Values[len(r.Values)-1].Value

		currentBytes := int(val)
		thresholdBytes := int(am.Spec.MemThresholdPercent) * int(podLimitsBytes) / 100
		// Check if the Pod needs redeployment, else continue onto the next Pod.
		if !needsRedeploy(currentBytes, thresholdBytes) {
			glog.V(2).Infof("Pod is operating normally: "+
				"%s | currentBytes: %d | thresholdBytes: %d",
				pod.Name, currentBytes, thresholdBytes)
			continue
		}

		// Redeploy the Pod with AppMonitor settings, if a redeploy is not already
		// in progress.
		// This operates a vertical autoscaling of the Pod.
		if !redeployInProgress(&pod) {
			glog.V(2).Infof("-------------------------------------------------------")
			glog.V(2).Infof("Pod *needs* redeployment: "+
				"%s | currentBytes: %d | thresholdBytes: %d",
				pod.Name, currentBytes, thresholdBytes)
			glog.V(2).Infof("-------------------------------------------------------")

			err := amc.redeployPodWithAppMonitor(&pod, &am)
			if err != nil {
				glog.V(2).Infof("Failed to vertically autoscale Pod: %s | "+
					"Error autoscaling: %v", pod.Name, err)
				continue
			}
		} else {
			glog.V(2).Infof("Redeploy in Progress for Pod: %s", pod.Name)
		}
	}
}

func needsRedeploy(currentBytes, thresholdBytes int) bool {
	if currentBytes >= thresholdBytes {
		return true
	}
	return false
}

func redeployInProgress(pod *v1.Pod) bool {
	if _, exists := pod.Annotations[AppMonitorAnnotationRedeployInProgress]; !exists {
		return false
	}
	return true
}

func (amc *AppMonitorController) redeployPodWithAppMonitor(pod *v1.Pod, am *AppMonitor) error {
	// Annotate the Pod in the cluster
	pod.Annotations[AppMonitorAnnotationRedeployInProgress] = "true"
	_, err := amc.clientSet.CoreV1().Pods(amc.namespace).Update(pod)
	if err != nil {
		return err
	}
	glog.V(2).Infof("Pod has been annotated: %s | Annotations: %s", pod.Name, pod.ObjectMeta.Annotations)

	// Create a new Pod in the cluster
	newPod := amc.newPodFromPod(pod, am)
	_, err = amc.clientSet.CoreV1().Pods(amc.namespace).Create(newPod)
	if err != nil {
		return err
	}
	glog.V(2).Infof("-----------------------------------------------------------")
	glog.V(2).Infof("AppMonitor autoscaled Pod: '%s' to '%s'", pod.Name, newPod.Name)
	glog.V(2).Infof("-----------------------------------------------------------")

	// Terminate Pod that crossed the threshold
	err = amc.clientSet.CoreV1().Pods(amc.namespace).Delete(pod.Name, nil)
	if err != nil {
		return err
	}
	glog.V(2).Infof("-----------------------------------------------------------")
	glog.V(2).Infof("Pod has been terminated: %s", pod.Name)
	glog.V(2).Infof("-----------------------------------------------------------")

	return nil
}

// Creates a new Pod for redeployment, based on the Pod being monitored
func (amc *AppMonitorController) newPodFromPod(pod *v1.Pod, am *AppMonitor) *v1.Pod {
	// Copy spec of the first container in the Pod
	newContainer := pod.Spec.Containers[0]

	// Reset VolumeMounts
	newContainer.VolumeMounts = nil

	// Set new resource limits based on the AppMonitor MemMultiplier
	podLimits, _ := newContainer.Resources.Limits["memory"]
	podLimitsBytes, _ := podLimits.AsInt64()
	newLimitsBytes := podLimitsBytes * int64(am.Spec.MemMultiplier)

	podRequests, _ := newContainer.Resources.Requests["memory"]
	podRequestsBytes, _ := podRequests.AsInt64()
	newRequestsBytes := podRequestsBytes * int64(am.Spec.MemMultiplier)

	newContainer.Resources.Limits = v1.ResourceList{
		v1.ResourceMemory: *resource.NewQuantity(
			newLimitsBytes,
			resource.BinarySI),
	}
	newContainer.Resources.Requests = v1.ResourceList{
		v1.ResourceMemory: *resource.NewQuantity(
			newRequestsBytes,
			resource.BinarySI),
	}

	// Create and return new Pod with AppMonitor settings applied
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-autoscaled", pod.Name),
			Namespace: pod.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				newContainer,
			},
		},
	}
}

// Query Prometheus over a range of time
func queryPrometheus(client prometheus.API, query string) (
	model.Value, error) {

	now := time.Now()
	start := now.Add(-1 * time.Second)
	results, err := client.QueryRange(
		context.Background(),
		query,
		prometheus.Range{
			Start: start,
			End:   now,
			Step:  time.Second,
		},
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Extract the Values of a Prometheus query for results of
//type prometheus.Matrix
func getMatrixValuesFromResults(results model.Value) []*model.SampleStream {
	// Type assert the interface to a model.Matrix
	matrix := results.(model.Matrix)

	// Type convert the matrix to a model.SampleStream to extract its stream of
	// values holding the data
	ss := []*model.SampleStream(matrix)

	return ss
}
