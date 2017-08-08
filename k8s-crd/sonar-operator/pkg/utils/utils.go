package utils

import (
	"github.com/golang/glog"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Build kubeconfig for use with clients.
// The kubeconfig file can either be passed in as a param, or attempted to be
// retrieved from the in-cluster ServiceAccount.
func BuildKubeConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		glog.V(2).Infof("kubeconfig file: %s", kubeconfig)
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	glog.V(2).Info("kubeconfig file: using InClusterConfig.")
	return rest.InClusterConfig()
}

// Attempt to deep copy an empty interface into a Pod.
func CopyObjToPod(obj interface{}) (*v1.Pod, error) {
	objCopy, err := api.Scheme.Copy(obj.(*v1.Pod))
	if err != nil {
		return nil, err
	}

	pod := objCopy.(*v1.Pod)
	if pod.ObjectMeta.Annotations == nil {
		pod.ObjectMeta.Annotations = make(map[string]string)
	}
	return pod, nil
}

// Attempt to deep copy an empty interface into a PodList.
func CopyObjToPods(obj []interface{}) ([]v1.Pod, error) {
	pods := []v1.Pod{}

	for _, o := range obj {
		pod, err := CopyObjToPod(o)
		if err != nil {
			glog.Errorf("Failed to copy pod object for podList: %v", err)
			return nil, err
		}
		pods = append(pods, *pod)
	}

	return pods, nil
}

// Select only the Pods with the annotation from the PodList.
func SelectAnnotatedPods(pods []v1.Pod, annotation string) ([]v1.Pod, error) {
	annotatedPods := []v1.Pod{}
	for _, pod := range pods {
		if _, exists := pod.ObjectMeta.Annotations[annotation]; !exists {
			glog.V(2).Infof("Skipping selection of non-annotated Pod: %s", pod.Name)
			continue
		}
		annotatedPods = append(annotatedPods, pod)
	}
	return annotatedPods, nil
}

// Attempt to deep copy an empty interface into a ThirdPartyResource.
func CopyObjToThirdPartyResource(obj interface{}) (*v1beta1.ThirdPartyResource, error) {
	objCopy, err := api.Scheme.Copy(obj.(*v1beta1.ThirdPartyResource))
	if err != nil {
		return nil, err
	}

	tpr := objCopy.(*v1beta1.ThirdPartyResource)
	if tpr.Annotations == nil {
		tpr.Annotations = make(map[string]string)
	}
	return tpr, nil
}
