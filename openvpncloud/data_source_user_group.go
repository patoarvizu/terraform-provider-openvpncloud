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

func dataSourceUserGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupRead,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpn_region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internet_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_device": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	userGroup, err := c.GetUserGroup(d.Get("name").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if userGroup == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User group not found",
			Detail:   fmt.Sprintf("User group not found"),
		})
		return diags
	}
	d.Set("user_group_id", userGroup.Id)
	d.Set("name", userGroup.Name)
	d.Set("vpn_region_ids", userGroup.VpnRegionIds)
	d.Set("internet_access", userGroup.InternetAccess)
	d.Set("max_device", userGroup.MaxDevice)
	d.Set("system_subnets", userGroup.SystemSubnets)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
