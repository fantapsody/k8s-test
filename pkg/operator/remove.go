package operator

import (
	"github.com/fantapsody/k8s-test/pkg/util"
	"github.com/golang/glog"
	olmClientVersioned "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"
	apiV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restClient "k8s.io/client-go/rest"
)

const csvName = "zookeeper-operator.v0.0.1"

func Remove(kubeConfig *restClient.Config, namespace string) error {
	kubeClient, err := kubernetes.NewForConfig(util.GetConfigSafe())
	if err != nil {
		panic(err)
	}

	// make sure the namespace exists
	_, e := kubeClient.CoreV1().Namespaces().Get(namespace, apiV1.GetOptions{})
	if e != nil {
		panic(e)
	}

	oc, e := olmClientVersioned.NewForConfig(kubeConfig)
	if e != nil {
		panic(e)
	}

	// delete all subscriptions
	subList, e := oc.OperatorsV1alpha1().Subscriptions(namespace).List(apiV1.ListOptions{})
	if e != nil {
		return e
	}
	for _, item := range subList.Items {
		if e = oc.OperatorsV1alpha1().Subscriptions(namespace).Delete(item.Name, &apiV1.DeleteOptions{}); e != nil {
			glog.Errorf("Failed to remove subscription %s: %s", item.Name, e.Error())
		} else {
			glog.Infof("Removed subscription %s", item.Name)
		}
	}

	catList, e := oc.OperatorsV1alpha1().CatalogSources(namespace).List(apiV1.ListOptions{})
	if e != nil {
		return e
	}
	for _, item := range catList.Items {
		if e = oc.OperatorsV1alpha1().CatalogSources(namespace).Delete(item.Name, &apiV1.DeleteOptions{}); e != nil {
			glog.Errorf("Failed to remove catalog source %s: %s", item.Name, e.Error())
		} else {
			glog.Infof("Removed catalog source %s", item.Name)
		}
	}

	// delete all CSVs
	csvList, e := oc.OperatorsV1alpha1().ClusterServiceVersions(namespace).List(apiV1.ListOptions{})
	if e != nil {
		return e
	}
	for _, item := range csvList.Items {
		glog.Infof("To remove item %s", item.Name)
		if e = oc.OperatorsV1alpha1().ClusterServiceVersions(namespace).Delete(item.Name, &apiV1.DeleteOptions{}); e != nil {
			glog.Errorf("Failed to remove csv %s: %s", item.Name, e.Error())
		} else {
			glog.Infof("Removed csv %s", item.Name)
		}
	}

	if e = oc.OperatorsV1().OperatorGroups(namespace).Delete(operatorGroupName, &apiV1.DeleteOptions{}); e != nil {
		glog.Errorf("Failed to remove operator group %s: %s", operatorGroupName, e.Error())
	} else {
		glog.Infof("Removed operator group %s", operatorGroupName)
	}

	return nil
}
