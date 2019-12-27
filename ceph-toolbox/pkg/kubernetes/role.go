package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	rbacApiV1Beta1 "k8s.io/api/rbac/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func RoleHandler(config, body string) error {
	return buildRole(config, body)
}

func buildRole(config, body string) error {
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

	role := rbacApiV1Beta1.Role{}
	if err := yaml.Unmarshal([]byte(body), &role); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := rbac.Roles(role.Namespace).Get(role.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.Roles(role.Namespace).Create(&role); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
