module github.com/kcp-dev/generic-controlplane

go 1.22.2

replace (
	k8s.io/api => github.com/mjudeikis/kubernetes/staging/src/k8s.io/api v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/apiextensions-apiserver => github.com/mjudeikis/kubernetes/staging/src/k8s.io/apiextensions-apiserver v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/apimachinery => github.com/mjudeikis/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/apiserver => github.com/mjudeikis/kubernetes/staging/src/k8s.io/apiserver v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/client-go => github.com/mjudeikis/kubernetes/staging/src/k8s.io/client-go v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/cloud-provider => github.com/mjudeikis/kubernetes/staging/src/k8s.io/cloud-provider v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/cluster-bootstrap => github.com/mjudeikis/kubernetes/staging/src/k8s.io/cluster-bootstrap v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/component-base => github.com/mjudeikis/kubernetes/staging/src/k8s.io/component-base v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/component-helpers => github.com/mjudeikis/kubernetes/staging/src/k8s.io/component-helpers v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/controller-manager => github.com/mjudeikis/kubernetes/staging/src/k8s.io/controller-manager v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/csi-translation-lib => github.com/mjudeikis/kubernetes/staging/src/k8s.io/csi-translation-lib v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/dynamic-resource-allocation => github.com/mjudeikis/kubernetes/staging/src/k8s.io/dynamic-resource-allocation v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/kms => github.com/mjudeikis/kubernetes/staging/src/k8s.io/kms v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/kube-aggregator => github.com/mjudeikis/kubernetes/staging/src/k8s.io/kube-aggregator v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/kubelet => github.com/mjudeikis/kubernetes/staging/src/k8s.io/kubelet v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/kubernetes => github.com/mjudeikis/kubernetes v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/kubernetes/pkg/kubeapiserver => github.com/mjudeikis/kubernetes/staging/src/k8s.io/apiserver v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/mount-utils => github.com/mjudeikis/kubernetes/staging/src/k8s.io/mount-utils v0.0.0-20240512101111-3f9e84a36b95
	k8s.io/pod-security-admission => github.com/mjudeikis/kubernetes/staging/src/k8s.io/pod-security-admission v0.0.0-20240512101111-3f9e84a36b95
)

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.6.0
	github.com/kcp-dev/kcp v0.24.0
	github.com/kcp-dev/kcp/cli v0.24.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace
	k8s.io/apiextensions-apiserver v0.30.0
	k8s.io/apimachinery v0.30.0
	k8s.io/apiserver v0.30.0
	k8s.io/client-go v1.5.2
	k8s.io/component-base v0.30.0
	k8s.io/klog/v2 v2.120.1
	k8s.io/kube-aggregator v0.30.0
	k8s.io/kubernetes v1.30.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/evanphx/json-patch v5.9.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/cel-go v0.20.1 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.7.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runc v1.1.12 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.14.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.9 // indirect
	go.etcd.io/etcd/api/v3 v3.5.13 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.13 // indirect
	go.etcd.io/etcd/client/v2 v2.305.13 // indirect
	go.etcd.io/etcd/client/v3 v3.5.13 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.13 // indirect
	go.etcd.io/etcd/raft/v3 v3.5.13 // indirect
	go.etcd.io/etcd/server/v3 v3.5.13 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.51.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.51.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/oauth2 v0.20.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20240509183442-62759503f434 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240509183442-62759503f434 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240509183442-62759503f434 // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.30.0 // indirect
	k8s.io/cloud-provider v0.30.0 // indirect
	k8s.io/cluster-bootstrap v0.30.0 // indirect
	k8s.io/component-helpers v0.30.0 // indirect
	k8s.io/controller-manager v0.30.0 // indirect
	k8s.io/dynamic-resource-allocation v0.30.0 // indirect
	k8s.io/kms v0.30.0 // indirect
	k8s.io/kube-openapi v0.0.0-20240430033511-f0e62f92d13f // indirect
	k8s.io/kubelet v0.30.0 // indirect
	k8s.io/mount-utils v0.30.0 // indirect
	k8s.io/pod-security-admission v0.30.0 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.30.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
