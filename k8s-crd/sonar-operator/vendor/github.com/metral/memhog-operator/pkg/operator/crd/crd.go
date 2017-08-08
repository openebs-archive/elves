package crd

import (
	"fmt"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateCustomResourceDefinition(clientSet apiextensionsclient.Interface, name, domain, kind, resourceNamePlural, version string) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	crdName := fmt.Sprintf("%s.%s", resourceNamePlural, domain)

	// initialize the CRD if it does not exist
	crd, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			crd := &apiextensionsv1beta1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: crdName,
				},
				Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
					Group:   domain,
					Version: version,
					Scope:   apiextensionsv1beta1.NamespaceScoped,
					Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
						Kind:   kind,
						Plural: resourceNamePlural,
					},
				},
			}

			result, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
			if err != nil {
				return nil, err
			}

			fmt.Printf("CREATED: %#v\nFROM: %#v\n", result, crd)

			// wait for CRD being established
			err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
				crd, err = clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crdName, metav1.GetOptions{})
				if err != nil {
					return false, err
				}
				for _, cond := range crd.Status.Conditions {
					switch cond.Type {
					case apiextensionsv1beta1.Established:
						if cond.Status == apiextensionsv1beta1.ConditionTrue {
							return true, err
						}
					case apiextensionsv1beta1.NamesAccepted:
						if cond.Status == apiextensionsv1beta1.ConditionFalse {
							fmt.Printf("Name conflict: %v\n", cond.Reason)
						}
					}
				}
				return false, err
			})
			if err != nil {
				deleteErr := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crdName, nil)
				if deleteErr != nil {
					return nil, utilerrors.NewAggregate([]error{err, deleteErr})
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		fmt.Printf("SKIPPING: already exists %#v\n", crd)
	}

	return crd, nil
}
