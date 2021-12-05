package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceNetworkRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRoutesRead,
		Schema: map[string]*schema.Schema{
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkRoutesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routes, err := c.GetRoutes(d.Get("network_item_id").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	configRoutes := make([]map[string]interface{}, len(routes))
	for i, r := range routes {
		route := make(map[string]interface{})
		routeType := r.Type
		route["type"] = routeType
		if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
			route["value"] = r.Subnet
		} else if routeType == client.RouteTypeDomain {
			route["value"] = r.Domain
		}
		configRoutes[i] = route
	}
	d.Set("routes", configRoutes)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
