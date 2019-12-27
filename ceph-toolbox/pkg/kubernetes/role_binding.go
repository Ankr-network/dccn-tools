package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rbacApiV1Beta1 "k8s.io/api/rbac/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func RoleBindingHandler(config, body string) error {
	return buildRoleBinding(config, body)
}

func buildRoleBinding(config, body string) error {

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

	roleBinding := rbacApiV1Beta1.RoleBinding{}

	if err := yaml.Unmarshal([]byte(body), &roleBinding); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := rbac.RoleBindings(roleBinding.Namespace).Get(roleBinding.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.RoleBindings(roleBinding.Namespace).Create(&roleBinding); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
