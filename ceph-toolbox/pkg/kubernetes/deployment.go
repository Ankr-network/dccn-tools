package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apisv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func DeploymentHandler(config, body string) error {
	return buildDeployment(config, body)
}

func buildDeployment(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	dc, err := apisv1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}

	deploy := apiv1.Deployment{}

	if err := yaml.Unmarshal([]byte(body), &deploy); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := dc.Deployments(deploy.Namespace).Get(deploy.Name, v1.GetOptions{}); err != nil {
		if _, err := dc.Deployments(deploy.Namespace).Create(&deploy); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
