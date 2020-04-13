package yaml

import (
	"testing"
)

type Pvc struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   struct {
		Name      string
		Namespace string
	}
	Spec struct {
		AccessModes []string `yaml:"accessModes"`
		Resources   struct {
			Requests struct {
				Storage string
			}
		}
	}
}

func TestRun(t *testing.T) {
	pvc := Pvc{}
	err := LoadFromYaml("pvc.yaml", &pvc)
	if err != nil {
		t.Error(err)
	}
	str, err := StructToYaml(&pvc)
	if err != nil {
		t.Error(err)
	}
}
