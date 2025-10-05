package main

import (
	"os"

	"github.com/kptdev/krm-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/helmfn"
)

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(helmfn.Process)); err != nil {
		os.Exit(1)
	}
}
