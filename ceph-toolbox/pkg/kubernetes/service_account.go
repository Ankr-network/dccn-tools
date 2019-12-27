package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func ServiceAccountHandler(config, body string) error {
	return buildServiceAccount(config, body)
}

func buildServiceAccount(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	coreV1Client, err := corev1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	sa := apiv1.ServiceAccount{}
	if err := yaml.Unmarshal([]byte(body), &sa); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := coreV1Client.ServiceAccounts(sa.Namespace).Get(sa.Name, v1.GetOptions{}); err != nil {
		if _, err = coreV1Client.ServiceAccounts(sa.Namespace).Create(&sa); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
