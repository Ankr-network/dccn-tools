package ceph

import (
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	glog.Info(strings.Repeat("==", 10), "start to install  ceph cluster", strings.Repeat("==", 10))
	glog.Info("prepare to install cluster condition")
	// initialize environment
	InitEnv(cmd)
	glog.Info("prepare work over")
	glog.Info(strings.Repeat("==", 50))
	glog.Info("config basic environment")
	// config cluster common
	if err := ConfigClusterCommon(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("config work over")
	glog.Info(strings.Repeat("==", 50))
	glog.Info("install cluster operator")
	// install cluster operator
	if err := InstallClusterOperator(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("operator install work over")
	glog.Info(strings.Repeat("==", 50))
	glog.Info("install ceph storage cluster")
	// install cluster until it runs normally
	if err := InstallCluster(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("ceph storage install work over")
	glog.Info(strings.Repeat("==", 50))
	// initialize ceph storage pool
	glog.Info("initialize ceph storage pool")
	if err := InstallCephBlockPool(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("ceph storage pool initialize work over")
	glog.Info(strings.Repeat("==", 50))
	// set ceph as default storage device
	glog.Info("set ceph as default storage device")
	if err := InstallClusterStorageClass(cmd); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("set default storage device work over")

	// all work over, well done!
	glog.Info("congratulations, the cluster have been installed!")
}
