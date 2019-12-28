package ceph

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func Uninstall(cmd *cobra.Command, args []string) {

	// delete sc
	if err := DelClusterStorageClass(cmd); err != nil {
		glog.Error(err)
		return
	}

	// delete storage pool
	if err := DelCephBlockPool(cmd); err != nil {
		glog.Error(err)
		return
	}

	// delete cluster
	if err := DeleteCluster(cmd); err != nil {
		glog.Error(err)
		return
	}

	// delete operator
	if err := DeleteClusterOperator(cmd); err != nil {
		glog.Error(err)
		return
	}

	// delete common
	if err := DelClusterComm(cmd); err != nil {
		glog.Error(err)
		return
	}
}
