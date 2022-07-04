package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"connectors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_item_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_item_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_region_id": {
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

func dataSourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	host, err := c.GetHostByName(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("name", host.Name)
	d.Set("internet_access", host.InternetAccess)
	d.Set("system_subnets", host.SystemSubnets)
	d.Set("connectors", getConnectorsSlice(&host.Connectors))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
