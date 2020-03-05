package operator

import (
	"github.com/fantapsody/k8s-test/pkg/util"
	"github.com/golang/glog"
	olmV1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	olmClientVersioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apiV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restClient "k8s.io/client-go/rest"
	"time"
)

const operatorhubCatalogSource = "operatorhubio-catalog"
const operatorGroupName = "pulsar-operatorgroup"
const pulsarCatalogSrcName = "pulsar-manifests"
const pulsarOperatorsImageName = "gcr.io/affable-ray-226821/streamnative/operator-manifest/pulsar-operators:latest"

const zookeeperCatalogSrcName = "zookeeper-manifests"

const zookeeperManifestsImageName = "gcr.io/affable-ray-226821/streamnative/operator-manifest/zookeeper-operator:snapshot-latest"

//const zookeeperManifestsImageName = "quay.io/fantapsody/operator-index-pulsar-manifests:latest"
const zookeeperSubscriptionName = "zookeeper-operator-subscription"
const zookeeperPackageName = "zookeeper-operator"

const bookkeeperCatalogSrcName = "bookkeeper-manifests"
const bookkeeperManifestsImageName = "gcr.io/affable-ray-226821/streamnative/operator-manifest/bookkeeper-operator:snapshot-latest"
const bookkeeperSubscriptionName = "bookkeeper-operator-subscription"
const bookkeeperPackageName = "bookkeeper-operator"

const prometheusSubscriptionName = "prometheus-operator-subscription"
const prometheusPackageName = "prometheus"

func Install(kubeConfig *restClient.Config, namespace string) error {
	kubeClient, err := kubernetes.NewForConfig(util.GetConfigSafe())
	if err != nil {
		return err
	}
	if err := EnsureEnvironment(kubeClient, namespace); err != nil {
		return err
	}
	if err := InstallOperator(kubeClient, kubeConfig, namespace); err != nil {
		return err
	}
	return nil
}

func EnsureEnvironment(kubeClient *kubernetes.Clientset, namespace string) error {
	namespaceClient := kubeClient.CoreV1().Namespaces()
	_, e := namespaceClient.Get(namespace, apiV1.GetOptions{})
	if e != nil {
		if errors.IsNotFound(e) {
			_, e = namespaceClient.Create(&coreV1.Namespace{
				ObjectMeta: apiV1.ObjectMeta{
					Name: namespace,
				},
			})
			if e != nil {
				return e
			}
			glog.Infof("Created namespace %s", namespace)
		} else {
			return e
		}
	} else {
		glog.Infof("Namespace %s exists", namespace)
	}
	return nil
}

func InstallOperator(kubeClient *kubernetes.Clientset, kubeConfig *restClient.Config, namespace string) error {
	oc, err := olmClientVersioned.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	if err := ensureCatalogSource(oc, namespace, pulsarCatalogSrcName, pulsarOperatorsImageName); err != nil {
		return err
	}
	//if err := ensureCatalogSource(oc, namespace, zookeeperCatalogSrcName, zookeeperManifestsImageName); err != nil {
	//	return err
	//}
	//if err := ensureCatalogSource(oc, namespace, bookkeeperCatalogSrcName, bookkeeperManifestsImageName); err != nil {
	//	return err
	//}
	if err := ensureOperatorGroup(oc, namespace, operatorGroupName); err != nil {
		return err
	}
	if err := ensureSubscription(oc, namespace, zookeeperSubscriptionName, pulsarCatalogSrcName, namespace, zookeeperPackageName, "alpha"); err != nil {
		return err
	}
	if err := ensureSubscription(oc, namespace, bookkeeperSubscriptionName, pulsarCatalogSrcName, namespace, bookkeeperPackageName, "alpha"); err != nil {
		return err
	}
	if err := ensureSubscription(oc, namespace, prometheusSubscriptionName, operatorhubCatalogSource, "olm", prometheusPackageName, "beta"); err != nil {
		return err
	}
	return nil
}

func ensureCatalogSource(oc *olmClientVersioned.Clientset, namespace, name, imageName string) error {
	catalogSource, e := oc.OperatorsV1alpha1().CatalogSources(namespace).Create(&v1alpha1.CatalogSource{
		ObjectMeta: apiV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType: "grpc",
			Image:      imageName,
			UpdateStrategy: &v1alpha1.UpdateStrategy{
				RegistryPoll: &v1alpha1.RegistryPoll{
					Interval: &apiV1.Duration{
						Duration: 1 * time.Minute,
					},
				},
			},
		},
	})
	if e != nil {
		if errors.IsAlreadyExists(e) {
			glog.Infof("Catalog source %s already exists", name)
		} else {
			return e
		}
	} else {
		glog.Infof("Created catalog source %s", catalogSource.Name)
	}
	return nil
}

func ensureOperatorGroup(oc *olmClientVersioned.Clientset, namespace, name string) error {
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
		if errors.IsAlreadyExists(e) {
			glog.Infof("Operator group %s already exists", name)
		} else {
			return e
		}
	} else {
		glog.Infof("Created operator group %s", operatorGroup.Name)
	}
	return nil
}

func ensureSubscription(oc *olmClientVersioned.Clientset, namespace, name, catalogSource, catalogSourceNamespace, packageName, channel string) error {
	subscription, e := oc.OperatorsV1alpha1().Subscriptions(namespace).Create(&v1alpha1.Subscription{
		ObjectMeta: apiV1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Package:                packageName,
			Channel:                channel,
			CatalogSourceNamespace: catalogSourceNamespace,
			CatalogSource:          catalogSource,
		},
	})
	if e != nil {
		if errors.IsAlreadyExists(e) {
			glog.Infof("Operator group %s already exists", name)
		} else {
			return e
		}
	} else {
		glog.Infof("Created operator group %s", subscription.Name)
	}
	return nil
}

func ensureCSV(oc *olmClientVersioned.Clientset, namespace string) {

}
