package veeam

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider ... provider for veeam
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP of the veeam server",
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_SERVER_IP", nil),
			},

			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server Port",
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_SERVER_PORT", nil),
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name of the veeam Server",
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_SERVER_USERNAME", nil),
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The user password of the VEEAM Server",
				DefaultFunc: schema.EnvDefaultFunc("VEEAM_SERVER_PASSWORD", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"veeam_job_vm": resourceVeeamJobVM(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ServerIP: d.Get("server_ip").(string),
		Port:     d.Get("port").(int),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	log.Printf("[DEBUG] Connecting to veeam backup server.......")
	return config, nil
}
