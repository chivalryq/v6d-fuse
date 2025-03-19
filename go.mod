module github.com/chivalryq/v6d-fuse

go 1.21

require github.com/spf13/pflag v1.0.5

require (
	github.com/apache/arrow/go/v11 v11.0.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.29.0 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/controller-runtime v0.17.2 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)

require (
	github.com/hanwen/go-fuse/v2 v2.7.2
	github.com/v6d-io/v6d/go/vineyard v0.0.0-20241202050933-fa7f6907b499
	golang.org/x/sys v0.18.0 // indirect
)

replace (
	github.com/go-logr/logr => github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr => github.com/go-logr/zapr v0.2.0
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.2.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.7.2
)

replace github.com/v6d-io/v6d/go/vineyard => ../v6d/go/vineyard
