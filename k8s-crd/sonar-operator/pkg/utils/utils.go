package utils

import (
	"fmt"

	"github.com/golang/glog"
        apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
        apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func HasCustomResourceDefinition(clientSet apiextensionsclient.Interface, name, domain, kind, resourceNamePlural, version string) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	crdName := fmt.Sprintf("%s.%s", resourceNamePlural, domain)

	crd, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	} 

	//TODO Verify the parameters
	return crd, nil
}
