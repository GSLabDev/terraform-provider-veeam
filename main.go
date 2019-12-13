package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/GSLabDev/terraform-provider-veeam/veeam"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: veeam.Provider})
}
