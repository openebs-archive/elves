package utils

import (
	"github.com/golang/glog"
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
