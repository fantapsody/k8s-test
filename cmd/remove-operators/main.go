package main

import (
	"github.com/fantapsody/k8s-test/pkg/operator"
	"github.com/fantapsody/k8s-test/pkg/util"
)

func main() {
	operator.Remove(util.GetConfigSafe(), "default")
}
