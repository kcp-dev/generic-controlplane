package batteries

import (
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/namespace/lifecycle"
	validatingadmissionpolicy "k8s.io/apiserver/pkg/admission/plugin/policy/validating"
	"k8s.io/apiserver/pkg/admission/plugin/resourcequota"
	mutatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	validatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
	controlplaneapiserver "k8s.io/kubernetes/pkg/controlplane/apiserver"
	"k8s.io/kubernetes/plugin/pkg/admission/admit"
	certapproval "k8s.io/kubernetes/plugin/pkg/admission/certificates/approval"
	"k8s.io/kubernetes/plugin/pkg/admission/certificates/ctbattest"

	certsigning "k8s.io/kubernetes/plugin/pkg/admission/certificates/signing"
	certsubjectrestriction "k8s.io/kubernetes/plugin/pkg/admission/certificates/subjectrestriction"
	"k8s.io/kubernetes/plugin/pkg/admission/defaulttolerationseconds"
	"k8s.io/kubernetes/plugin/pkg/admission/deny"
	"k8s.io/kubernetes/plugin/pkg/admission/eventratelimit"
	"k8s.io/kubernetes/plugin/pkg/admission/gc"
	"k8s.io/kubernetes/plugin/pkg/admission/namespace/autoprovision"
	"k8s.io/kubernetes/plugin/pkg/admission/namespace/exists"
	"k8s.io/kubernetes/plugin/pkg/admission/serviceaccount"
	// Admission policies
)

type Batteries struct {
	list BatteriesList

	BatteriesArgs []string
}

type Battery string

type BatteriesList map[Battery]BatterySpec

type BatterySpec struct {
	// Enabled indicates whether the battery is enabled.
	Enabled bool

	// GroupNames is the list of group names that the battery is responsible for.
	// If disabled, the battery will not be registered for these groups.
	GroupNames []string
}

const (
	// BatteryLeases is the name of the lease battery.
	BatteryLeases Battery = "leases"
	// BatteryAuthentication is the name of the authentication battery.
	BatteryAuthentication Battery = "authentication"
	// BatteryAuthorization is the name of the authorization battery.
	BatteryAuthorization Battery = "authorization"
	// BatteryAdmission is the name of the admission battery.
	BatteryAdmission Battery = "admission"
	// BatteryFlowControl is the name of the flow control battery.
	BatteryFlowControl Battery = "flowcontrol"
	// BatteryCRDs is the name of the CRD battery.
	BatteryCRDs Battery = "crds"
)

var (
	// The generic features.
	defaultBatteries = map[Battery]BatterySpec{
		BatteryLeases:         {Enabled: false, GroupNames: []string{"coordination.k8s.io"}},
		BatteryAuthentication: {Enabled: false, GroupNames: []string{"authentication.k8s.io", "rbac.authentication.k8s.io"}},
		BatteryAuthorization:  {Enabled: false, GroupNames: []string{"authorization.k8s.io", "rbac.authorization.k8s.io"}},
		BatteryAdmission:      {Enabled: false, GroupNames: []string{"admissionregistration.k8s.io"}},
		BatteryFlowControl:    {Enabled: false, GroupNames: []string{"flowcontrol.apiserver.k8s.io"}},
		BatteryCRDs:           {Enabled: false, GroupNames: []string{"apiextensions.k8s.io"}},
	}
)

func (b Battery) String() string {
	return string(b)
}

func New() Batteries {
	b := Batteries{
		list: make(BatteriesList, len(defaultBatteries)),
	}
	for name, spec := range defaultBatteries {
		b.list[name] = spec
	}
	return b
}

func (b Batteries) Enable(name Battery) {
	_b := b.list[name]
	_b.Enabled = true
	b.list[name] = _b
}

func (b Batteries) Disable(name Battery) {
	_b := b.list[name]
	_b.Enabled = false
	b.list[name] = _b
}

func (b Batteries) IsEnabled(name Battery) bool {
	spec, ok := b.list[name]
	return ok && spec.Enabled
}

// RegisterAllAdmissionPlugins registers all admission plugins based on the batteries configuration.
func (b Batteries) RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	admit.Register(plugins) // DEPRECATED as no real meaning
	autoprovision.Register(plugins)
	lifecycle.Register(plugins)
	exists.Register(plugins)
	serviceaccount.Register(plugins)
	eventratelimit.Register(plugins)
	gc.Register(plugins)
	certapproval.Register(plugins)
	certsigning.Register(plugins)
	ctbattest.Register(plugins)
	certsubjectrestriction.Register(plugins)
	mutatingwebhook.Register(plugins)
	validatingadmissionpolicy.Register(plugins)
	validatingwebhook.Register(plugins)
	resourcequota.Register(plugins)
	deny.Register(plugins)
}

func (b Batteries) DefaultOffAdmissionPlugins() sets.Set[string] {
	defaultOnPlugins := sets.New(
		lifecycle.PluginName, // NamespaceLifecycle
		// limitranger.PluginName,               // LimitRanger
		serviceaccount.PluginName,           // ServiceAccount
		resourcequota.PluginName,            // ResourceQuota
		certapproval.PluginName,             // CertificateApproval
		certsigning.PluginName,              // CertificateSigning
		ctbattest.PluginName,                // ClusterTrustBundleAttest
		certsubjectrestriction.PluginName,   // CertificateSubjectRestriction
		defaulttolerationseconds.PluginName, // DefaultTolerationSeconds
	)

	if b.IsEnabled(BatteryAdmission) {
		defaultOnPlugins.Insert(
			mutatingwebhook.PluginName,           // MutatingAdmissionWebhook
			validatingwebhook.PluginName,         // ValidatingAdmissionWebhook
			validatingadmissionpolicy.PluginName, // ValidatingAdmissionPolicy, only active when feature gate ValidatingAdmissionPolicy is enabled
		)
	}

	return sets.New[string](AllOrderedPlugins...).Difference(defaultOnPlugins)
}

func (b Batteries) containsAndDisabled(name string) bool {
	for _, spec := range b.list {
		if slices.Contains(spec.GroupNames, name) && !spec.Enabled {
			return true
		}
	}
	return false
}

func (b Batteries) FilterStorageProviders(input []controlplaneapiserver.RESTStorageProvider) []controlplaneapiserver.RESTStorageProvider {
	var result []controlplaneapiserver.RESTStorageProvider
	for _, rest := range input {
		if b.containsAndDisabled(rest.GroupName()) {
			continue
		}
		result = append(result, rest)
	}
	return result
}
