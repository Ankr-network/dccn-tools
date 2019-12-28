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

func RoleHandler(config, body string) error {
	return buildRole(config, body)
}

func buildRole(config, body string) error {
	rbac, role, e := roleComm(config, body)
	if e != nil {
		return e
	}

	if _, err := rbac.Roles(role.Namespace).Get(role.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.Roles(role.Namespace).Create(&role); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func roleComm(config string, body string) (*rbacv1beta1.RbacV1beta1Client, rbacApiV1Beta1.Role, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.Role{}, err
	}
	rbac, err := rbacv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.Role{}, err
	}
	role := rbacApiV1Beta1.Role{}
	if err := yaml.Unmarshal([]byte(body), &role); err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.Role{}, err
	}
	return rbac, role, nil
}

func DelRoleHandler(config, body string) error {
	return delRole(config, body)
}

func delRole(config, body string) error {
	rbac, role, e := roleComm(config, body)
	if e != nil {
		return e
	}

	if _, err := rbac.Roles(role.Namespace).Get(role.Name, v1.GetOptions{}); err == nil {
		if err := rbac.Roles(role.Namespace).Delete(role.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := rbac.Roles(role.Namespace).Get(role.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
