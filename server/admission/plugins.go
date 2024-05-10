package admission

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/namespace/lifecycle"
	mutatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	validatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/validating"

	// Admission policies
	"k8s.io/kubernetes/plugin/pkg/admission/admit"
	"k8s.io/kubernetes/plugin/pkg/admission/defaulttolerationseconds"
	"k8s.io/kubernetes/plugin/pkg/admission/serviceaccount"
)

// AllOrderedPlugins is the list of all the plugins in order.
var AllOrderedPlugins = []string{
	admit.PluginName, // AlwaysAdmit
	//autoprovision.PluginName, // NamespaceAutoProvision
	////lifecycle.PluginName,     // NamespaceLifecycle
	//exists.PluginName,        // NamespaceExists
	//limitranger.PluginName,            // LimitRanger
	serviceaccount.PluginName, // ServiceAccount
	//eventratelimit.PluginName, // EventRateLimit
	//gc.PluginName,             // OwnerReferencesPermissionEnforcement
	//certapproval.PluginName,           // CertificateApproval
	//certsigning.PluginName,            // CertificateSigning
	//ctbattest.PluginName,              // ClusterTrustBundleAttest
	//certsubjectrestriction.PluginName, // CertificateSubjectRestriction

	// new admission plugins should generally be inserted above here
	// webhook, resourcequota, and deny plugins must go at the end
	//mutatingwebhook.PluginName,           // MutatingAdmissionWebhook
	///validatingadmissionpolicy.PluginName, // ValidatingAdmissionPolicy
	//validatingwebhook.PluginName,         // ValidatingAdmissionWebhook
	//resourcequota.PluginName,             // ResourceQuota
	//deny.PluginName, // AlwaysDeny
}

// RegisterAllAdmissionPlugins registers all admission plugins.
// The order of registration is irrelevant, see AllOrderedPlugins for execution order.
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	admit.Register(plugins) // DEPRECATED as no real meaning
	//autoprovision.Register(plugins)
	//lifecycle.Register(plugins)
	//exists.Register(plugins)
	//limitranger.Register(plugins)
	serviceaccount.Register(plugins)
	//eventratelimit.Register(plugins)
	//gc.Register(plugins)
	//certapproval.Register(plugins)
	////certsigning.Register(plugins)
	//ctbattest.Register(plugins)
	//certsubjectrestriction.Register(plugins)
	////mutatingwebhook.Register(plugins)
	//validatingadmissionpolicy.Register(plugins)
	//validatingwebhook.Register(plugins)
	//resourcequota.Register(plugins)
	//deny.Register(plugins)
}

// DefaultOffAdmissionPlugins get admission plugins off by default for kube-apiserver.
func DefaultOffAdmissionPlugins() sets.Set[string] {
	defaultOnPlugins := sets.New(
		lifecycle.PluginName, // NamespaceLifecycle
		//	limitranger.PluginName,               // LimitRanger
		serviceaccount.PluginName,           // ServiceAccount
		defaulttolerationseconds.PluginName, // DefaultTolerationSeconds
		mutatingwebhook.PluginName,          // MutatingAdmissionWebhook
		validatingwebhook.PluginName,        // ValidatingAdmissionWebhook
		//	resourcequota.PluginName,             // ResourceQuota
		//certapproval.PluginName,              // CertificateApproval
		//certsigning.PluginName,               // CertificateSigning
		//ctbattest.PluginName,                 // ClusterTrustBundleAttest
		//certsubjectrestriction.PluginName,    // CertificateSubjectRestriction
		//validatingadmissionpolicy.PluginName, // ValidatingAdmissionPolicy, only active when feature gate ValidatingAdmissionPolicy is enabled
	)

	return sets.New[string](AllOrderedPlugins...).Difference(defaultOnPlugins)
}
