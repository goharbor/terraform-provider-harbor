package main

import (
	"bitbucket.org/bestsellerit/terraform-provider-harbor/provider"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider})
}
