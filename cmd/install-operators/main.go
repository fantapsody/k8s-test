package main

import (
	"github.com/fantapsody/k8s-test/pkg/operator"
	"github.com/fantapsody/k8s-test/pkg/util"
)

func main() {
	if err := operator.Install(util.GetConfigSafe(), "default"); err != nil {
		panic(err)
	}
}
