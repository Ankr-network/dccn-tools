package kubernetes

import (
	"time"

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
	coreV1Client, sa, e := serviceAccountComm(config, body)
	if e != nil {
		return e
	}

	if _, err := coreV1Client.ServiceAccounts(sa.Namespace).Get(sa.Name, v1.GetOptions{}); err != nil {
		if _, err = coreV1Client.ServiceAccounts(sa.Namespace).Create(&sa); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func serviceAccountComm(config string, body string) (*corev1.CoreV1Client, apiv1.ServiceAccount, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.ServiceAccount{}, err
	}
	coreV1Client, err := corev1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.ServiceAccount{}, err
	}
	sa := apiv1.ServiceAccount{}
	if err := yaml.Unmarshal([]byte(body), &sa); err != nil {
		glog.Error(err)
		return nil, apiv1.ServiceAccount{}, err
	}
	return coreV1Client, sa, nil
}

func DelServiceAccountHandler(config, body string) error {
	return delServiceAccount(config, body)
}

func delServiceAccount(config, body string) error {
	coreV1Client, sa, e := serviceAccountComm(config, body)
	if e != nil {
		return e
	}
	if _, err := coreV1Client.ServiceAccounts(sa.Namespace).Get(sa.Name, v1.GetOptions{}); err == nil {
		if err = coreV1Client.ServiceAccounts(sa.Namespace).Delete(sa.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := coreV1Client.ServiceAccounts(sa.Namespace).Get(sa.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	return nil
}
