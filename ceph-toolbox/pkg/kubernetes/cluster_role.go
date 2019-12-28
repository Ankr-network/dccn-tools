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

func ClusterRoleHandler(config, body string) error {
	return buildClusterRole(config, body)
}

func buildClusterRole(config, body string) error {
	rbac, clusterRole, e := clusterRoleComm(config, body)
	if e != nil {
		return e
	}
	if _, err := rbac.ClusterRoles().Get(clusterRole.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.ClusterRoles().Create(&clusterRole); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func clusterRoleComm(config string, body string) (*rbacv1beta1.RbacV1beta1Client, rbacApiV1Beta1.ClusterRole, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.ClusterRole{}, err
	}
	rbac, err := rbacv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.ClusterRole{}, err
	}
	clusterRole := rbacApiV1Beta1.ClusterRole{}
	if err := yaml.Unmarshal([]byte(body), &clusterRole); err != nil {
		glog.Error(err)
		return nil, rbacApiV1Beta1.ClusterRole{}, err
	}
	return rbac, clusterRole, nil
}

func DelClusterRoleHandler(config, body string) error {
	return delClusterRole(config, body)
}

func delClusterRole(config, body string) error {
	rbac, clusterRole, e := clusterRoleComm(config, body)
	if e != nil {
		return e
	}
	if _, err := rbac.ClusterRoles().Get(clusterRole.Name, v1.GetOptions{}); err == nil {
		if err := rbac.ClusterRoles().Delete(clusterRole.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}

	// wait the component end
	for {
		if _, err := rbac.ClusterRoles().Get(clusterRole.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
