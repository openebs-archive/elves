package main

import (
	"os"

	"github.com/openebs/elves/k8s-crd/sonar-operator/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
