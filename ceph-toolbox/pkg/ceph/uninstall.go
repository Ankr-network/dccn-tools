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
	glog.Info("clear storage class over")

	// delete storage pool
	if err := DelCephBlockPool(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("clear ceph block pool over")

	// delete cluster
	if err := DeleteCluster(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("clear ceph cluster over.")

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
