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

type Battery string

type BatteriesList map[Battery]BatterySpec

type BatterySpec struct {
	// Enabled indicates whether the battery is enabled.
	Enabled bool

	// Description is a human-readable description of the battery.
	Description string

	// Groups is the list of group names that the battery is responsible for.
	// If disabled, the battery will not be registered for these groups.
	Groups []string
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
		BatteryLeases: {
			Enabled:     false,
			Groups:      []string{"coordination.k8s.io"},
			Description: "Leases are used to coordinate some operations between Kubernetes components"},
		BatteryAuthentication: {
			Enabled:     false,
			Groups:      []string{"authentication.k8s.io"},
			Description: "Authentication verifies the identity of the user",
		},
		BatteryAuthorization: {
			Enabled:     false,
			Groups:      []string{"authorization.k8s.io", "rbac.authorization.k8s.io"},
			Description: "Authorization decides whether a request is allowed",
		},
		BatteryAdmission: {
			Enabled:     false,
			Groups:      []string{"admissionregistration.k8s.io"},
			Description: "Admission controllers validate and mutate requests",
		},
		BatteryFlowControl: {
			Enabled:     false,
			Groups:      []string{"flowcontrol.apiserver.k8s.io"},
			Description: "Flow control limits number of requests processed at a time",
		},
		BatteryCRDs: {
			Enabled:     false,
			Groups:      []string{"apiextensions.k8s.io"},
			Description: "CustomResourceDefinitions (CRDs) allow definition of custom resources",
		},
	}
)

func (b Battery) String() string {
	return string(b)
}

func New() Options {
	b := Options{
		batteries: make(BatteriesList, len(defaultBatteries)),
	}
	for name, spec := range defaultBatteries {
		b.batteries[name] = spec
	}
	return b
}

func (b Options) Enable(name Battery) {
	_b := b.batteries[name]
	_b.Enabled = true
	b.batteries[name] = _b
}

func (b Options) Disable(name Battery) {
	_b := b.batteries[name]
	_b.Enabled = false
	b.batteries[name] = _b
}

func (b CompletedOptions) IsEnabled(name Battery) bool {
	spec, ok := b.batteries[name]
	return ok && spec.Enabled
}

// RegisterAllAdmissionPlugins registers all admission plugins based on the batteries configuration.
func (b CompletedOptions) RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
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

func (b CompletedOptions) DefaultOffAdmissionPlugins() sets.Set[string] {
	defaultOnPlugins := sets.New[string](
		lifecycle.PluginName, // NamespaceLifecycle
		// limitranger.PluginName,           // LimitRanger
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

func (b CompletedOptions) containsAndDisabled(name string) bool {
	for _, spec := range b.batteries {
		if slices.Contains(spec.Groups, name) && !spec.Enabled {
			return true
		}
	}
	return false
}

func (b CompletedOptions) FilterStorageProviders(input []controlplaneapiserver.RESTStorageProvider) []controlplaneapiserver.RESTStorageProvider {
	var result []controlplaneapiserver.RESTStorageProvider
	for _, rest := range input {
		if b.containsAndDisabled(rest.GroupName()) {
			continue
		}
		result = append(result, rest)
	}
	return result
}
