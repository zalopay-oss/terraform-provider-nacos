package nacos

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nacos "gitlab.zalopay.vn/top/cicd/terraform-provider-nacos/pkg/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NACOS_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NACOS_PASSWORD", nil),
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NACOS_ADDRESS", nil),
			},
			"context_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NACOS_CONTEXT_PATH", "nacos"),
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"nacos_configuration": resourceConfiguration(),
		},
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	c, err := nacos.NewClient(&nacos.Config{
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		Address:     d.Get("address").(string),
		ContextPath: d.Get("context_path").(string),
	})

	if err != nil {
		return nil, diag.Errorf("create nacos client error: %v", err)
	}
	return c, diags
}
