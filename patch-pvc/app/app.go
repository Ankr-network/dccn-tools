package app

// App about app-pod
type App struct {
	namespace string
	appName   string
	pvc       string
	pv        string
}

// GetPvcAndPv  kubectl get pvc -n <namespace> get pvc and pv
func (a *App) GetPvcAndPv() {}

// DownLoadPvc kubectl get pvc <pvc> -n <namespace> -o yaml > <pvc>.yaml  download pvc.yaml
func (a *App) DownLoadPvc() {}

// HandPvc handle pvc.yaml
func (a *App) HandPvc() {}

// DelPvc kubectl -n <namespace> delete pvc <pvc>
func (a *App) DelPvc() {}

// DelPv kubectl -n <namespace> delete pv <pv>
func (a *App) DelPv() {}

// PatchPvc kubectl -n <namespace> patch pvc <pvc> -p '{"metadata":{"finalizers": []}}' --type=merge
func (a *App) PatchPvc() {}

//PatchPv kubectl patch pv <pv> -p '{"metadata":{"finalizers": []}}' --type=merge
func (a *App) PatchPv() {}

// ApplyPvc kubectl apply -f <pvc>.yaml
func (a *App) ApplyPvc() {}

// CreateApp  create app struct
func CreateApp(namespace, appName string) App {
	return App{
		namespace,
		appName,
		"",
		"",
	}
}
