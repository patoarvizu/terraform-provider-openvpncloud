package openvpncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Managed by Terraform",
				ValidateFunc: validation.StringLenBetween(1, 120),
			},
			"internet_access": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      client.InternetAccessLocal,
				ValidateFunc: validation.StringInSlice([]string{client.InternetAccessBlocked, client.InternetAccessGlobalInternet, client.InternetAccessLocal}, false),
			},
			"system_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connector": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpn_region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_item_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_item_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_v4_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_v6_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	var connectors []client.Connector
	configConnectors := d.Get("connector").([]interface{})
	for _, c := range configConnectors {
		connectors = append(connectors, client.Connector{
			Name:        c.(map[string]interface{})["name"].(string),
			VpnRegionId: c.(map[string]interface{})["vpn_region_id"].(string),
		})
	}
	h := client.Host{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		InternetAccess: d.Get("internet_access").(string),
		Connectors:     connectors,
	}
	host, err := c.CreateHost(h)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(host.Id)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The connector for this host needs to be set up manually",
		Detail:   "Terraform only creates the OpenVPN Cloud connector object for this host, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	host, err := c.GetHostById(d.Id())
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if host == nil {
		d.SetId("")
		return diags
	}
	d.Set("name", host.Name)
	d.Set("description", host.Description)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)
	if len(d.Get("connector").([]interface{})) > 0 {
		configConnector := d.Get("connector").([]interface{})[0].(map[string]interface{})
		connectorName := configConnector["name"].(string)
		err = d.Set("connector", getConnectorSlice(host.Connectors, host.Id, connectorName))
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("connector") {
		old, new := d.GetChange("connector")
		oldSlice := old.([]interface{})
		newSlice := new.([]interface{})
		if len(oldSlice) == 0 && len(newSlice) == 1 {
			// This happens when importing the resource
			newConnector := client.Connector{
				Name:            newSlice[0].(map[string]interface{})["name"].(string),
				VpnRegionId:     newSlice[0].(map[string]interface{})["vpn_region_id"].(string),
				NetworkItemType: client.NetworkItemTypeHost,
			}
			_, err := c.AddConnector(newConnector, d.Id())
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}
		} else {
			oldMap := oldSlice[0].(map[string]interface{})
			newMap := newSlice[0].(map[string]interface{})
			if oldMap["name"].(string) != newMap["name"].(string) || oldMap["vpn_region_id"].(string) != newMap["vpn_region_id"].(string) {
				newConnector := client.Connector{
					Name:            newMap["name"].(string),
					VpnRegionId:     newMap["vpn_region_id"].(string),
					NetworkItemType: client.NetworkItemTypeHost,
				}
				_, err := c.AddConnector(newConnector, d.Id())
				if err != nil {
					return append(diags, diag.FromErr(err)...)
				}
				if len(oldMap["id"].(string)) > 0 {
					// This can sometimes happen when importing the resource
					err = c.DeleteConnector(oldMap["id"].(string), d.Id(), oldMap["network_item_type"].(string))
					if err != nil {
						return append(diags, diag.FromErr(err)...)
					}
				}
			}
		}
	}
	if d.HasChanges("name", "description", "internet_access") {
		_, newName := d.GetChange("name")
		_, newDescription := d.GetChange("description")
		_, newAccess := d.GetChange("internet_access")
		err := c.UpdateHost(client.Host{
			Id:             d.Id(),
			Name:           newName.(string),
			Description:    newDescription.(string),
			InternetAccess: newAccess.(string),
		})
		if err != nil {
			return append(diags, diag.FromErr(err)...)
		}
	}
	return append(diags, resourceHostRead(ctx, d, m)...)
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	hostId := d.Id()
	err := c.DeleteHost(hostId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
