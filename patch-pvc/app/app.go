package app

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Ankr-network/dccn-tools/patch-pvc/app/pvc"
	"github.com/Ankr-network/dccn-tools/patch-pvc/funs"
	"github.com/Ankr-network/dccn-tools/patch-pvc/lib/yaml"
)

// App about app-pod
type App struct {
	namespace string
	appName   string
	pvc       string
	pv        string
}

// GetPvcAndPv  kubectl get pvc -n <namespace> get pvc and pv
func (a *App) GetPvcAndPv() error {
	out, err := funs.Exec("kubectl", "get", "pvc", "-n", a.namespace)
	if err != nil {
		return errors.New(fmt.Sprintf("kubectl get pvc -n %s  err:%v", a.namespace, err))
	}
	bf := bufio.NewReader(strings.NewReader(out))
	bf.ReadString('\n')
	str, _ := bf.ReadString('\n')
	strs := strings.Fields(str)
	a.pvc = strs[0]
	a.pv = strs[2]
	return nil
}

// DownLoadPvc kubectl get pvc <pvc> -n <namespace> -o yaml > <pvc>.yaml  download pvc.yaml
func (a *App) DownLoadPvc() error {
	out, err := funs.Exec("kubectl", "get", "pvc", a.pvc, "-n", a.namespace, "-o", "yaml")
	if err != nil {
		return errors.New(fmt.Sprintf("kubectl get pvc -n %s err:%v", a.namespace, err))
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", a.pvc), []byte(out), 0744)
	if err != nil {
		return errors.New(fmt.Sprintf("DownLoadPvc WriteFile err:%s", err))
	}
	return nil
}

// HandPvc handle pvc.yaml
func (a *App) HandPvc(pvc interface{}) error {

	err := yaml.LoadFromYaml(fmt.Sprintf("%s.yaml", a.pvc), pvc)
	if err != nil {
		return errors.New(fmt.Sprintf("HandPvc LoadFromYaml err:%s", err))
	}
	str, err := yaml.StructToYaml(pvc)
	if err != nil {
		return errors.New(fmt.Sprintf("HandPvc StructToYaml err:%s", err))
	}
	err = ioutil.WriteFile("pvc.yaml", []byte(str), 0744)
	if err != nil {
		return errors.New(fmt.Sprintf("HandPvc WriteFile err:%s", err))
	}
	return nil
}

// DelPvc kubectl -n <namespace> delete pvc <pvc>
func (a *App) DelPvc() {
	funs.Exec("kubectl", "-n", a.namespace, "delete", "pvc", a.pvc)
}

// DelPv kubectl -n <namespace> delete pv <pv>
func (a *App) DelPv() {
	funs.Exec("kubectl", "-n", a.namespace, "delete", "pv", a.pv)
}

// PatchPvc kubectl -n <namespace> patch pvc <pvc> -p '{"metadata":{"finalizers": []}}' --type=merge
func (a *App) PatchPvc() error {
	_, err := funs.Exec("kubectl", "-n", a.namespace, "patch", "pvc", a.pvc, "-p", `{"metadata":{"finalizers": []}}`, "--type=merge")
	if err != nil {
		return errors.New(fmt.Sprintf(`kubectl -n %s patch pvc %s -p '{"metadata":{"finalizers": []}}' --type=merge`, a.namespace, a.pvc))
	}
	return nil
}

//PatchPv kubectl patch pv <pv> -p '{"metadata":{"finalizers": []}}' --type=merge
func (a *App) PatchPv() error {
	_, err := funs.Exec("kubectl", "patch", "pv", a.pv, "-p", `{"metadata":{"finalizers": []}}`, "--type=merge")
	if err != nil {
		return errors.New(fmt.Sprintf(`kubectl patch pv %s -p '{"metadata":{"finalizers": []}}' --type=merge`, a.pv))
	}
	return nil
}

// ApplyPvc kubectl apply -f <pvc>.yaml
func (a *App) ApplyPvc() error {
	_, err := funs.Exec("kubectl", "apply", "-f", "pvc.yaml")
	if err != nil {
		return errors.New(fmt.Sprintf("kubectl apply -f pvc.yaml err:%s", err))
	}
	return nil
}

// Run 执行流程
func (a *App) Run() {
	a.GetPvcAndPv()
	a.DownLoadPvc()
	pvc := pvc.Pvc1{}
	a.HandPvc(&pvc)
	a.DelPvc()
	a.DelPv()
	a.PatchPvc()
	a.PatchPv()
	a.ApplyPvc()
}

// CreateApp  create app struct
func CreateApp(namespace, appName string) App {
	return App{
		namespace,
		appName,
		"",
		"",
	}
}
