package kubernetes

import (
	"time"

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

	rbac, roleBinding, e := roleBindingComm(config, body)
	if e != nil {
		return e
	}

	if _, err := rbac.RoleBindings(roleBinding.Namespace).Get(roleBinding.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.RoleBindings(roleBinding.Namespace).Create(&roleBinding); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func roleBindingComm(config string, body string) (*rbacv1beta1.RbacV1beta1Client, rbacApiV1Beta1.RoleBinding, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.RoleBinding{}, err
	}
	rbac, err := rbacv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.RoleBinding{}, err
	}
	roleBinding := rbacApiV1Beta1.RoleBinding{}
	if err := yaml.Unmarshal([]byte(body), &roleBinding); err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.RoleBinding{}, err
	}
	return rbac, roleBinding, nil
}

func DelRoleBindingHandler(config, body string) error {
	return delRoleBinding(config, body)
}

func delRoleBinding(config, body string) error {

	rbac, roleBinding, e := roleBindingComm(config, body)
	if e != nil {
		return e
	}
	if _, err := rbac.RoleBindings(roleBinding.Namespace).Get(roleBinding.Name, v1.GetOptions{}); err == nil {
		if err := rbac.RoleBindings(roleBinding.Namespace).Delete(roleBinding.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := rbac.RoleBindings(roleBinding.Namespace).Get(roleBinding.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	return nil
}
