package yaml

import (
	"fmt"
	"io/ioutil"

	"errors"
	"gopkg.in/yaml.v2"
)
// LoadFromYaml yaml.conf load
func LoadFromYaml(path string, out interface{})error{
	yaml_byte,err := ioutil.ReadFile(path)
	if err != nil{
		return errors.New(fmt.Sprintf("yaml.LoadFormYaml ioutil.ReadFile err:%s",err))
	}
	
	err = yaml.Unmarshal(yaml_byte, out)
	if err != nil {
		return errors.New(fmt.Sprintf("yaml.LoadFormYaml Unmarshal err:%s",err))
	}
	return nil
}

// StructToYaml  struct to yaml string
func StructToYaml(out interface{})(string,error){
	d,err := yaml.Marshal(out)
	if err != nil {
		return "", errors.New(fmt.Sprintf("yaml.StructToYaml Marshal err:%s",err))
	}
	return string(d),nil
}