package app

import (
	"github.com/metral/memhog-operator/pkg/cmd"
	k8slogsutil "k8s.io/kubernetes/pkg/util/logs"
)

// Run the memhog-operator command
func Run() error {
	// Init logging
	k8slogsutil.InitLogs()
	defer k8slogsutil.FlushLogs()

	// Create & execute new command
	cmd, err := cmd.NewCmdMemHogOperator()
	if err != nil {
		return err
	}

	return cmd.Execute()
}
