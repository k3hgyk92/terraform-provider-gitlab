package main

import (
	"flag"

	"github.com/gitlabhq/terraform-provider-gitlab/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// A bogus change here.

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version), Debug: debugMode, ProviderAddr: "registry.terraform.io/providers/gitlabhq/gitlab"}
	plugin.Serve(opts)
}
