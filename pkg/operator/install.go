package operator

import (
	"github.com/fantapsody/k8s-test/pkg/util"
	"github.com/golang/glog"
	olmV1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	olmClientVersioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	coreV1 "k8s.io/api/core/v1"
	apiV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restClient "k8s.io/client-go/rest"
	"strings"
)

const catalogSrcName = "pulsar-manifests"
const imageName = "gcr.io/affable-ray-226821/streamnative/operator-index/pulsar:v0.0.1"
const operatorGroupName = "pulsar-operatorgroup"
const zookeeperSubscriptionName = "zookeeper-operator-subscription"
const zookeeperPackageName = "zookeeper-operator"

func Install(kubeConfig *restClient.Config, namespace string) {
	kubeClient, err := kubernetes.NewForConfig(util.GetConfigSafe())
	if err != nil {
		panic(err)
	}
	EnsureEnvironment(kubeClient, namespace)
	InstallOperator(kubeClient, kubeConfig, namespace)
}

func EnsureEnvironment(kubeClient *kubernetes.Clientset, namespace string) {
	namespaceClient := kubeClient.CoreV1().Namespaces()
	_, e := namespaceClient.Get(namespace, apiV1.GetOptions{})
	if e != nil {
		if strings.Contains(e.Error(), "not found") {
			_, e = namespaceClient.Create(&coreV1.Namespace{
				ObjectMeta: apiV1.ObjectMeta{
					Name: namespace,
				},
			})
			if e != nil {
				panic(e)
			}
			glog.Infof("Created namespace %s", namespace)
		} else {
			panic(e)
		}
	} else {
		glog.Infof("Namespace %s exists", namespace)
	}
}

func InstallOperator(kubeClient *kubernetes.Clientset, kubeConfig *restClient.Config, namespace string) {
	oc, err := olmClientVersioned.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	ensureCatalogSource(oc, namespace, catalogSrcName)
	ensureOperatorGroup(oc, namespace, operatorGroupName)
	ensureSubscription(oc, namespace, zookeeperSubscriptionName, catalogSrcName, zookeeperPackageName)

}

func ensureCatalogSource(oc *olmClientVersioned.Clientset, namespace, name string) {
	catalogSource, e := oc.OperatorsV1alpha1().CatalogSources(namespace).Create(&v1alpha1.CatalogSource{
		ObjectMeta: apiV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType: "grpc",
			Image:      imageName,
		},
	})
	if e != nil {
		if strings.Contains(e.Error(), "already exists") {
			glog.Infof("Catalog source %s already exists", name)
		} else {
			panic(e)
		}
	} else {
		glog.Infof("Created catalog source %s", catalogSource.Name)
	}
}

func ensureOperatorGroup(oc *olmClientVersioned.Clientset, namespace, name string) {
	operatorGroup, e := oc.OperatorsV1().OperatorGroups(namespace).Create(&olmV1.OperatorGroup{
		ObjectMeta: apiV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: olmV1.OperatorGroupSpec{
			TargetNamespaces: []string{namespace},
		},
	})
	if e != nil {
		if strings.Contains(e.Error(), "already exists") {
			glog.Infof("Operator group %s already exists", name)
		} else {
			panic(e)
		}
	} else {
		glog.Infof("Created operator group %s", operatorGroup.Name)
	}
}

func ensureSubscription(oc *olmClientVersioned.Clientset, namespace, name, catalogSource, packageName string) {
	subscription, e := oc.OperatorsV1alpha1().Subscriptions(namespace).Create(&v1alpha1.Subscription{
		ObjectMeta: apiV1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Package:                packageName,
			Channel:                "alpha",
			CatalogSourceNamespace: namespace,
			CatalogSource:          catalogSource,
		},
	})
	if e != nil {
		if strings.Contains(e.Error(), "already exists") {
			glog.Infof("Operator group %s already exists", name)
		} else {
			panic(e)
		}
	} else {
		glog.Infof("Created operator group %s", subscription.Name)
	}
}