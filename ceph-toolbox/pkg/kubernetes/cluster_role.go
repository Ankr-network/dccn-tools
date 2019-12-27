package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rbacApiV1Beta1 "k8s.io/api/rbac/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func ClusterRoleHandler(config, body string) error {
	return buildClusterRole(config, body)
}

func buildClusterRole(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	rbac, err := rbacv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	clusterRole := rbacApiV1Beta1.ClusterRole{}
	if err := yaml.Unmarshal([]byte(body), &clusterRole); err != nil {
		glog.Error(err)
		return err
	}
	if _, err := rbac.ClusterRoles().Get(clusterRole.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.ClusterRoles().Create(&clusterRole); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
