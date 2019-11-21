package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/nttcom/terraform-provider-fic/fic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: fic.Provider})
}
