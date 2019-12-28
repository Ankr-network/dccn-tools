package kubernetes

var (
	ApiHandler    map[string]func(config, body string) error
	DelApiHandler map[string]func(config, body string) error
)

func init() {
	ApiHandler = make(map[string]func(config, body string) error)
	ApiHandler[Namespace] = NamespaceHandler
	ApiHandler[CRD] = CustomResourceDefinitionHandler
	ApiHandler[ClusterRoleBinding] = ClusterRoleBindingHandler
	ApiHandler[ClusterRole] = ClusterRoleHandler
	ApiHandler[ServiceAccount] = ServiceAccountHandler
	ApiHandler[Role] = RoleHandler
	ApiHandler[RoleBinding] = RoleBindingHandler
	ApiHandler[PodSecurityPolicy] = PodSecurityPolicyHandler
	ApiHandler[Deployment] = DeploymentHandler

	DelApiHandler = make(map[string]func(config, body string) error)
	DelApiHandler[Namespace] = DelNamespaceHandler
	DelApiHandler[CRD] = DelCustomResourceDefinitionHandler
	DelApiHandler[ClusterRoleBinding] = DelClusterRoleBindingHandler
	DelApiHandler[ClusterRole] = DelClusterRoleHandler
	DelApiHandler[ServiceAccount] = DelServiceAccountHandler
	DelApiHandler[Role] = DelRoleHandler
	DelApiHandler[RoleBinding] = DelRoleBindingHandler
	DelApiHandler[PodSecurityPolicy] = DelPodSecurityPolicyHandler
	DelApiHandler[Deployment] = DelDeploymentHandler
}

const (
	Namespace          = "Namespace"
	CRD                = "CustomResourceDefinition"
	ClusterRoleBinding = "ClusterRoleBinding"
	ClusterRole        = "ClusterRole"
	Role               = "Role"
	ServiceAccount     = "ServiceAccount"
	RoleBinding        = "RoleBinding"
	PodSecurityPolicy  = "PodSecurityPolicy"
	Deployment         = "Deployment"
)
