package ceph

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rookApi "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	rookv1 "github.com/rook/rook/pkg/client/clientset/versioned/typed/ceph.rook.io/v1"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func InstallCephBlockPool(cmd *cobra.Command) error {
	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}

	cephBlockPool := rookApi.CephBlockPool{}
	if err := yaml.Unmarshal([]byte(ClusterReplicatedPool), &cephBlockPool); err != nil {
		glog.Error(err)
		return err
	}

	cephV1Client, err := rookv1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	if _, err := cephV1Client.CephClusters(cephBlockPool.Namespace).Get(cephBlockPool.Name, v1.GetOptions{}); err != nil {
		if _, err := cephV1Client.CephBlockPools(cephBlockPool.Namespace).Create(&cephBlockPool); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
