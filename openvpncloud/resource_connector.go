package openvpncloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectorCreate,
		ReadContext:   resourceConnectorRead,
		DeleteContext: resourceConnectorDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_item_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{client.NetworkItemTypeHost, client.NetworkItemTypeNetwork}, false),
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	}
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	networkItemId := d.Get("network_item_id").(string)
	networkItemType := d.Get("network_item_type").(string)
	vpnRegionId := d.Get("vpn_region_id").(string)
	connector := client.Connector{
		Name:            name,
		NetworkItemId:   networkItemId,
		NetworkItemType: networkItemType,
		VpnRegionId:     vpnRegionId,
	}
	conn, err := c.AddNetworkConnector(connector, networkItemId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(conn.Id)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Connector needs to be set up manually",
		Detail:   "Terraform only creates the OpenVPN Cloud connector object, but additional manual steps are required to associate a host in your infrastructure with this connector. Go to https://openvpn.net/cloud-docs/connector/ for more information.",
	})
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	connector, err := c.GetConnectorByName(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if connector == nil {
		d.SetId("")
		return append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Connector with name %s not found", d.Get("name").(string)),
		})
	} else {
		d.SetId(connector.Id)
		d.Set("name", connector.Name)
		d.Set("vpn_region_id", connector.VpnRegionId)
		d.Set("network_item_type", connector.NetworkItemType)
		d.Set("network_item_id", connector.NetworkItemId)
		d.Set("ip_v4_address", connector.IPv4Address)
		d.Set("ip_v6_address", connector.IPv6Address)
	}
	return diags
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	err := c.DeleteNetworkConnector(d.Id(), d.Get("network_item_id").(string), d.Get("network_item_type").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
