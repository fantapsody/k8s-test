module github.com/fantapsody/k8s-test

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20200130164400-12c06cfc05c4
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v8.0.0+incompatible
)

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm
	github.com/openshift/api => github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20190923180330-3b6373338c9b
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.0.0-beta.5.0.20200123114618-5e3c7d7eb86a

	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	k8s.io/apiserver => k8s.io/apiserver v0.17.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.3
	k8s.io/code-generator => k8s.io/code-generator v0.17.3
	k8s.io/component-base => k8s.io/component-base v0.17.3
	k8s.io/cri-api => k8s.io/cri-api v0.17.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.3
	k8s.io/kubectl => k8s.io/kubectl v0.17.3
	k8s.io/kubelet => k8s.io/kubelet v0.17.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.3
	k8s.io/metrics => k8s.io/metrics v0.17.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.3
)
