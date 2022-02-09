module github.com/replicatedhq/kubeflare

go 1.16

require (
	github.com/cloudflare/cloudflare-go v0.31.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.19.16
	k8s.io/apimachinery v0.19.16
	k8s.io/cli-runtime v0.19.16
	k8s.io/client-go v0.19.16
	sigs.k8s.io/controller-runtime v0.6.5
)

require (
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	k8s.io/code-generator v0.19.17-rc.0 // indirect
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	k8s.io/utils v0.0.0-20210802155522-efc7438f0176 // indirect
	sigs.k8s.io/controller-tools v0.4.1 // indirect
)

replace github.com/appscode/jsonpatch => github.com/gomodules/jsonpatch v2.0.1+incompatible
