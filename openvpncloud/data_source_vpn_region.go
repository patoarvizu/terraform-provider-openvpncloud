package openvpncloud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceVpnRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnRegionRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"continent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country_iso": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	vpnRegion, err := c.GetVpnRegion(d.Get("region_id").(string))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error requesting VPN regions",
			Detail:   fmt.Sprintf("Error requesting VPN regions %v", err),
		})
		return diags
	}
	d.Set("region_id", vpnRegion.Id)
	d.Set("continent", vpnRegion.Continent)
	d.Set("country", vpnRegion.Country)
	d.Set("country_iso", vpnRegion.CountryISO)
	d.Set("region_name", vpnRegion.RegionName)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
