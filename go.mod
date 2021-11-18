module github.com/replicatedhq/kubeflare

go 1.16

require (
	github.com/cloudflare/cloudflare-go v0.13.8
	github.com/evanphx/json-patch v4.9.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/onsi/gomega v1.10.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/procfs v0.0.11 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/tools v0.0.0-20200616195046-dc31b401abb5 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6 // indirect
	k8s.io/apimachinery v0.18.6
	k8s.io/cli-runtime v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
)

replace github.com/appscode/jsonpatch => github.com/gomodules/jsonpatch v2.0.1+incompatible
