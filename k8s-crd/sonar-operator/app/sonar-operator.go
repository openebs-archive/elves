package app

import (
	"github.com/openebs/elves/k8s-crd/sonar-operator/pkg/cmd"
	k8slogsutil "k8s.io/kubernetes/pkg/util/logs"
)

// Run the sonar-operator command
func Run() error {
	// Init logging
	k8slogsutil.InitLogs()
	defer k8slogsutil.FlushLogs()

	// Create & execute new command
	cmd, err := cmd.NewCmdSonarOperator()
	if err != nil {
		return err
	}

	return cmd.Execute()
}
