package kubernetes

import (
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	policyApiV1Beta1 "k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	policyv1beta1 "k8s.io/client-go/kubernetes/typed/policy/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func PodSecurityPolicyHandler(config, body string) error {
	return buildPodSecurityPolicy(config, body)
}

func buildPodSecurityPolicy(config, body string) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return err
	}
	exc, err := policyv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return err
	}
	podSecPolicy := policyApiV1Beta1.PodSecurityPolicy{}

	if err := yaml.Unmarshal([]byte(body), &podSecPolicy); err != nil {
		glog.Error(err)
		return err
	}

	if _, err := exc.PodSecurityPolicies().Get(podSecPolicy.Name, v1.GetOptions{}); err != nil {
		if _, err := exc.PodSecurityPolicies().Create(&podSecPolicy); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}
