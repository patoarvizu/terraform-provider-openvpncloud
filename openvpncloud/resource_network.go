package openvpncloud

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"egress": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"internet_access": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{client.InternetAccessBlocked, client.InternetAccessGlobalInternet, client.InternetAccessLocal}, false),
				Optional:     true,
			},
			"system_subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_route": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      client.RouteTypeIPV4,
							ValidateFunc: validation.StringInSlice([]string{client.RouteTypeIPV4, client.RouteTypeIPV6, client.RouteTypeDomain}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"default_connector": {
				Type:     schema.TypeSet,
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

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	configConnector := d.Get("default_connector").(*schema.Set)
	connectors := []client.Connector{
		{
			Name:        configConnector.List()[0].(map[string]interface{})["name"].(string),
			VpnRegionId: configConnector.List()[0].(map[string]interface{})["vpn_region_id"].(string),
		},
	}
	n := client.Network{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Egress:         d.Get("egress").(bool),
		InternetAccess: d.Get("internet_access").(string),
		Connectors:     connectors,
	}
	network, err := c.CreateNetwork(n)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(network.Id)
	configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
	defaultRoute, err := c.CreateRoute(network.Id, client.Route{
		Type:  configRoute["type"].(string),
		Value: configRoute["value"].(string),
	})
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defaultRouteWithIdSlice := make([]map[string]interface{}, 1)
	defaultRouteWithIdSlice[0] = map[string]interface{}{
		"id": defaultRoute.Id,
	}
	d.Set("default_route", defaultRouteWithIdSlice)
	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	network, err := c.GetNetwork(d.Get("name").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error requesting network",
			Detail:   fmt.Sprintf("Error requesting network %v", err),
		})
		return diags
	}
	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("egress", network.Egress)
	d.Set("internet_access", network.InternetAccess)
	d.Set("system_subnets", network.SystemSubnets)
	connector := d.Get("default_connector").(*schema.Set).List()[0].(map[string]interface{})
	connectorName := connector["name"].(string)
	conns := getSingleNetworkConnector(c, network.Id, connectorName)
	d.Set("default_connector", conns)
	configRoute := d.Get("default_route").([]interface{})[0].(map[string]interface{})
	route, err := c.GetRoute(d.Id(), configRoute["id"].(string))
	defaultRoute := []map[string]interface{}{
		{
			"id":   configRoute["id"].(string),
			"type": route.Type,
		},
	}
	if route.Type == client.RouteTypeIPV4 || route.Type == client.RouteTypeIPV6 {
		defaultRoute[0]["value"] = route.Subnet
	} else if route.Type == client.RouteTypeDomain {
		defaultRoute[0]["value"] = route.Domain
	}
	err = d.Set("default_route", defaultRoute)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("default_connector") {
		o, n := d.GetChange("default_connector")
		old := o.(*schema.Set).List()[0].(map[string]interface{})
		new := n.(*schema.Set).List()[0].(map[string]interface{})
		if old["name"].(string) != new["name"].(string) || old["vpn_region_id"].(string) != new["vpn_region_id"].(string) {
			newConnector := client.Connector{
				Name:        new["name"].(string),
				VpnRegionId: new["vpn_region_id"].(string),
			}
			_, err := c.AddNetworkConnector(newConnector, d.Id())
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Error adding connector %v", err),
				})
			}
			err = c.RemoveNetworkConnector(old["id"].(string), d.Id())
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Error removing connector: %v", err),
				})
			}
		}
	}
	if d.HasChange("default_route") {
		o, n := d.GetChange("default_route")
		old := o.([]interface{})[0].(map[string]interface{})
		new := n.([]interface{})[0].(map[string]interface{})
		networkId := d.Id()
		routeId := old["id"]
		routeType := new["type"]
		routeValue := new["value"]
		route := client.Route{
			Id:    routeId.(string),
			Type:  routeType.(string),
			Value: routeValue.(string),
		}
		err := c.UpdateRoute(networkId, route)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return append(diags, resourceNetworkRead(ctx, d, m)...)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Deleting an openvpncloud_network resource is not supported.",
		Detail:   "Deleting a network is not supported by the OpenVPN cloud API yet. This operation only removed it from the Terraform state, but you'll need to manually delete it from the web console.",
	})
	return diags
}

func getSingleNetworkConnector(c *client.Client, networkId string, connectorName string) *schema.Set {
	connectorsSet := schema.NewSet(connectorsHash, []interface{}{})
	networkConnectors, _ := c.GetConnectorsForNetwork(networkId)
	for _, c := range networkConnectors {
		if c.NetworkItemId == networkId && c.Name == connectorName {
			connector := make(map[string]interface{})
			connector["id"] = c.Id
			connector["name"] = c.Name
			connector["network_item_id"] = c.NetworkItemId
			connector["network_item_type"] = c.NetworkItemType
			connector["vpn_region_id"] = c.VpnRegionId
			connector["ip_v4_address"] = c.IPv4Address
			connector["ip_v6_address"] = c.IPv6Address
			connectorsSet.Add(connector)
			break
		}
	}
	return connectorsSet
}

func connectorsHash(v interface{}) int {
	m, ok := v.(map[string]interface{})

	if !ok {
		return 0
	}

	if v, ok := m["name"].(string); ok {
		return stringHashcode(v)
	}

	return 0
}

func stringHashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}
