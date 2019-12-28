package ceph

import (
	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/kubernetes"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func InstallClusterStorageClass(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}
	if err := kubernetes.StorageClassHandler(config, ClusterStorageClass); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func DelClusterStorageClass(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}
	if err := kubernetes.DelStorageClassHandler(config, ClusterStorageClass); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
