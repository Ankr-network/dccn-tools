package ceph

import (
	"errors"
	"time"

	"github.com/Ankr-network/dccn-tools/ceph-toolbox/pkg/kubernetes"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apisv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func InstallClusterOperator(cmd *cobra.Command) error {

	var m map[string]interface{}

	if err := yaml.Unmarshal([]byte(ClusterOperator), &m); err != nil {
		glog.Error(err)
		return err
	}

	handleKey, ok := m[Kind].(string)
	if !ok {
		return errors.New("handle key assert failed")
	}

	config, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		glog.Error(err)
		return err
	}

	if err := kubernetes.ApiHandler[handleKey](config, ClusterOperator); err != nil {
		glog.Error(err)
		return err
	}

	// wait to return until the operator to is running
	if err := waitOperatorToRunning(config, ClusterOperator); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func waitOperatorToRunning(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	dc, err := apisv1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	deploy := apiv1.Deployment{}

	if err := yaml.Unmarshal([]byte(body), &deploy); err != nil {
		glog.Error(err)
		return err
	}

	// get replicate set associate deployment
	for {
		if dp, err := dc.Deployments(deploy.Namespace).Get(deploy.Name, metav1.GetOptions{}); err != nil {
			glog.Error(err)
			return err
		} else {
			if dp.Status.AvailableReplicas == dp.Status.ReadyReplicas && dp.Status.AvailableReplicas == dp.Status.Replicas {
				break
			}
		}
		time.Sleep(time.Second)
	}
	return nil
}
