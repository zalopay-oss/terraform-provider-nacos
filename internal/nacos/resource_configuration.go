package nacos

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	nacos "gitlab.zalopay.vn/top/cicd/terraform-provider-nacos/pkg/client"
)

func resourceConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		CreateContext: resourceConfigurationCreate,
		ReadContext:   resourceConfigurationRead,
		UpdateContext: resourceConfigurationUpdate,
		DeleteContext: resourceConfigurationDelete,
	}
}

func resourceConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nacos.Client)

	configuration := &nacos.Configuration{
		Namespace:   d.Get("namespace").(string),
		Group:       d.Get("group").(string),
		Key:         d.Get("key").(string),
		Value:       d.Get("value").(string),
		Description: d.Get("description").(string),
	}
	err := client.PublishConfiguration(ctx, configuration)
	if err != nil {
		return diag.Errorf("failed to create configuration = %+v: %v", *configuration, err)
	}

	d.SetId(convToResourceId(configuration.Namespace, configuration.Group, configuration.Key))

	return resourceConfigurationRead(ctx, d, meta)
}

func resourceConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nacos.Client)

	configurationId, err := convToConfigurationId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	configuration, err := client.GetConfiguration(ctx, configurationId)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"namespace":   configuration.Namespace,
		"group":       configuration.Group,
		"key":         configuration.Key,
		"value":       configuration.Value,
		"description": configuration.Description,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(convToResourceId(configuration.Namespace, configuration.Group, configuration.Key))

	return nil
}

func resourceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nacos.Client)
	if d.HasChanges("value", "description") {
		configuration := &nacos.Configuration{
			Namespace:   d.Get("namespace").(string),
			Group:       d.Get("group").(string),
			Key:         d.Get("key").(string),
			Value:       d.Get("value").(string),
			Description: d.Get("description").(string),
		}
		err := client.PublishConfiguration(ctx, configuration)
		if err != nil {
			return diag.Errorf("failed to update configuration = %+v: %v", *configuration, err)
		}
	}

	return resourceConfigurationRead(ctx, d, meta)
}

func resourceConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*nacos.Client)
	configurationId, err := convToConfigurationId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.DeleteConfiguration(ctx, configurationId)
	if err != nil {
		return diag.FromErr(err)
	}
	if !res {
		return diag.Errorf("failed to delete configurationId: %+v", *configurationId)
	}

	return nil
}
