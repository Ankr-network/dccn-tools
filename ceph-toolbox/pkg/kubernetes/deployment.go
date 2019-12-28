package kubernetes

import (
	"time"

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
	dc, deploy, e := deployComm(config, body)
	if e != nil {
		return e
	}

	if _, err := dc.Deployments(deploy.Namespace).Get(deploy.Name, v1.GetOptions{}); err != nil {
		if _, err := dc.Deployments(deploy.Namespace).Create(&deploy); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func deployComm(config string, body string) (*apisv1.AppsV1Client, apiv1.Deployment, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.Deployment{}, err
	}
	dc, err := apisv1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, apiv1.Deployment{}, err
	}
	deploy := apiv1.Deployment{}
	if err := yaml.Unmarshal([]byte(body), &deploy); err != nil {
		glog.Error(err)
		return nil, apiv1.Deployment{}, err
	}
	return dc, deploy, nil
}

func DelDeploymentHandler(config, body string) error {
	return delDeployment(config, body)
}

func delDeployment(config, body string) error {
	dc, deploy, e := deployComm(config, body)
	if e != nil {
		return e
	}
	if _, err := dc.Deployments(deploy.Namespace).Get(deploy.Name, v1.GetOptions{}); err == nil {
		if err := dc.Deployments(deploy.Namespace).Delete(deploy.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := dc.Deployments(deploy.Namespace).Get(deploy.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	return nil
}
