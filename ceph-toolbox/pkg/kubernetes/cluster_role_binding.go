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

func ClusterRoleBindingHandler(config, body string) error {
	return buildClusterRoleBinding(config, body)
}

func buildClusterRoleBinding(config, body string) error {
	err, rbac, clusterRoleBinding, e := clusterRoleBindingComm(config, body)
	if e != nil {
		return e
	}

	if _, err = rbac.ClusterRoleBindings().Get(clusterRoleBinding.Name, v1.GetOptions{}); err != nil {
		if _, err := rbac.ClusterRoleBindings().Create(&clusterRoleBinding); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func clusterRoleBindingComm(config string, body string) (error, *rbacv1beta1.RbacV1beta1Client, rbacApiV1Beta1.ClusterRoleBinding, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, nil, rbacApiV1Beta1.ClusterRoleBinding{}, err
	}
	rbac, err := rbacv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, nil, rbacApiV1Beta1.ClusterRoleBinding{}, err
	}
	clusterRoleBinding := rbacApiV1Beta1.ClusterRoleBinding{}
	if err := yaml.Unmarshal([]byte(body), &clusterRoleBinding); err != nil {
		glog.Error(err)
		return nil, nil, rbacApiV1Beta1.ClusterRoleBinding{}, err
	}
	return err, rbac, clusterRoleBinding, nil
}

func DelClusterRoleBindingHandler(config, body string) error {
	return delClusterRoleBinding(config, body)
}

func delClusterRoleBinding(config, body string) error {
	err, rbac, clusterRoleBinding, e := clusterRoleBindingComm(config, body)
	if e != nil {
		return e
	}

	if _, err = rbac.ClusterRoleBindings().Get(clusterRoleBinding.Name, v1.GetOptions{}); err == nil {
		if err := rbac.ClusterRoleBindings().Delete(clusterRoleBinding.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err = rbac.ClusterRoleBindings().Get(clusterRoleBinding.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return nil
}
