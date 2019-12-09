package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/terraform-providers/terraform-provider-veeam/veeam"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: veeam.Provider})
}
