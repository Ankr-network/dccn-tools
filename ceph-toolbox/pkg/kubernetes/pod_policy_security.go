package kubernetes

import (
	"time"

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
	exc, podSecPolicy, e := podSecComm(config, body)
	if e != nil {
		return e
	}

	if _, err := exc.PodSecurityPolicies().Get(podSecPolicy.Name, v1.GetOptions{}); err != nil {
		if _, err := exc.PodSecurityPolicies().Create(&podSecPolicy); err != nil {
			glog.Error(err)
			return err
		}
	}

	return nil
}

func podSecComm(config string, body string) (*policyv1beta1.PolicyV1beta1Client, policyApiV1Beta1.PodSecurityPolicy, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		glog.Error(err)
		return nil, policyApiV1Beta1.PodSecurityPolicy{}, err
	}
	exc, err := policyv1beta1.NewForConfig(cfg)
	if err != nil {
		glog.Error(err)
		return nil, policyApiV1Beta1.PodSecurityPolicy{}, err
	}
	podSecPolicy := policyApiV1Beta1.PodSecurityPolicy{}
	if err := yaml.Unmarshal([]byte(body), &podSecPolicy); err != nil {
		glog.Error(err)
		return nil, policyApiV1Beta1.PodSecurityPolicy{}, err
	}
	return exc, podSecPolicy, nil
}

func DelPodSecurityPolicyHandler(config, body string) error {
	return delPodSecurityPolicy(config, body)
}

func delPodSecurityPolicy(config, body string) error {
	exc, podSecPolicy, e := podSecComm(config, body)
	if e != nil {
		return e
	}
	if _, err := exc.PodSecurityPolicies().Get(podSecPolicy.Name, v1.GetOptions{}); err == nil {
		if err := exc.PodSecurityPolicies().Delete(podSecPolicy.Name, &v1.DeleteOptions{}); err != nil {
			glog.Error(err)
			return err
		}
	}
	// wait the component end
	for {
		if _, err := exc.PodSecurityPolicies().Get(podSecPolicy.Name, v1.GetOptions{}); err == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	return nil
}
