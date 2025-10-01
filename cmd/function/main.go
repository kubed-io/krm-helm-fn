package main

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/kubed-io/krm-helm-fn/pkg/helmfn"
)

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(helmfn.Process)); err != nil {
		os.Exit(1)
	}
}