package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"devices": {
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
						"description": {
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

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	user, err := c.GetUser(d.Get("username").(string), d.Get("role").(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.Set("user_id", user.Id)
	d.Set("username", user.Username)
	d.Set("role", user.Role)
	d.Set("email", user.Email)
	d.Set("auth_type", user.AuthType)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("group_id", user.GroupId)
	d.Set("status", user.Status)
	d.Set("devices", getUserDevicesSlice(&user.Devices))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func getUserDevicesSlice(userDevices *[]client.Device) []interface{} {
	devices := make([]interface{}, len(*userDevices), len(*userDevices))
	for i, d := range *userDevices {
		device := make(map[string]interface{})
		device["id"] = d.Id
		device["name"] = d.Name
		device["description"] = d.Description
		device["ip_v4_address"] = d.IPv4Address
		device["ip_v6_address"] = d.IPv6Address
		devices[i] = device
	}
	return devices
}
