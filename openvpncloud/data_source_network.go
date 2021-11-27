package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"egress": {
				Type:     schema.TypeBool,
				Computed: true,
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
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	network, err := c.GetNetworkByName(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("network_id", network.Id)
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	d.Set("routes", getRoutesSlice(&network.Routes))
	d.Set("connectors", getNetworkConnectorsSlice(&network.Connectors))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func getRoutesSlice(networkRoutes *[]client.Route) []interface{} {
	routes := make([]interface{}, len(*networkRoutes), len(*networkRoutes))
	for i, r := range *networkRoutes {
		route := make(map[string]interface{})
		route["id"] = r.Id
		route["subnet"] = r.Subnet
		route["type"] = r.Type
		routes[i] = route
	}
	return routes
}

func getNetworkConnectorsSlice(networkConnectors *[]client.Connector) []interface{} {
	connectors := make([]interface{}, len(*networkConnectors), len(*networkConnectors))
	for i, c := range *networkConnectors {
		connector := make(map[string]interface{})
		connector["id"] = c.Id
		connector["name"] = c.Name
		connector["network_item_id"] = c.NetworkItemId
		connector["network_item_type"] = c.NetworkItemType
		connector["vpn_region_id"] = c.VpnRegionId
		connector["ip_v4_address"] = c.IPv4Address
		connector["ip_v6_address"] = c.IPv6Address
		connectors[i] = connector
	}
	return connectors
}
