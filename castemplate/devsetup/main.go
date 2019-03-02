package main

import (
	"fmt"

	"github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/maya/pkg/client/k8s"
	"github.com/openebs/maya/pkg/engine"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testEngine struct {
	engine        engine.Interface
	defaultConfig []v1alpha1.Config
	openebsConfig []v1alpha1.Config
}

func main() {
	newK8sClient, err := k8s.NewK8sClient("")
	if err != nil {
		fmt.Println("error in getting clientset")
		fmt.Println(err)
		return
	}

	key := "demo-template"
	cast, err := newK8sClient.GetOEV1alpha1CAST(key, metav1.GetOptions{})
	if err != nil {
		fmt.Println("error in getting cas template")
		fmt.Println(err)
	}
	engine, err := engine.New(
		cast,
		key,
		map[string]interface{}{},
	)
	if err != nil {
		fmt.Println("error in creating machine")
		fmt.Println(err)
		return
	}

	data, err := engine.Run()
	if err != nil {
		fmt.Println("error in executing machine")
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}
