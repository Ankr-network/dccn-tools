package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rbacApiV1Beta1 "k8s.io/api/rbac/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func ClusterRoleBindingHandler(config, body string) error {
	return buildClusterRoleBinding(config, body)
}

func buildClusterRoleBinding(config, body string) error {
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

	clusterRoleBinding := rbacApiV1Beta1.ClusterRoleBinding{}

	if err := yaml.Unmarshal([]byte(body), &clusterRoleBinding); err != nil {
		glog.Error(err)
		return err
	}

	if _, err = rbac.ClusterRoleBindings().Get(clusterRoleBinding.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.ClusterRoleBindings().Create(&clusterRoleBinding); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
