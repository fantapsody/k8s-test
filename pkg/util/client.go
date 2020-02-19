package util

import (
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetConfigString() string {
	return filepath.Join(homedir.HomeDir(), ".kube", "config")
}

func GetConfig() (*restclient.Config, error) {
	return clientcmd.BuildConfigFromFlags("", GetConfigString())
}

func GetConfigSafe() *restclient.Config {
	config, e := GetConfig()
	if e != nil {
		panic(e)
	}
	return config
}
