package operator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// #############################################################################

// Note: The following code is custom settings particular to the new CRD in the
// cluster.

// #############################################################################

const CRDName string = "submarines"
const ResourceName string = "submarine"
const ResourceNamePlural string = "submarines"
const Domain string = "kubedemo.com"
const Version string = "v1alpha1"

// An Submarine redeploys an app when resource limits are exceeded.
type SubmarineSpec struct {
	nation string `json:nation`
}

// #############################################################################

// Note: The following code is boilerplate code needed to satisfy the
// AppMontior as a resource in the cluster in terms of how it expects CRD's to
// be created, operate and used.

// #############################################################################

type Submarine struct {
	// TODO: add Name field for Submarine as its currently missing
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`

	Spec SubmarineSpec `json:"spec"`
}

type SubmarineList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`

	Items []Submarine `json:"items"`
}
