package openvpncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 120),
			},
			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
			},
			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 20),
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"devices": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
						},
						"description": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 120),
						},
						"ipv4_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	username := d.Get("username").(string)
	// OpenVPN Cloud API has a bug that does not allow setting the role during the user's creation
	role := "MEMBER"
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	groupId := d.Get("group_id").(string)
	configDevices := d.Get("devices").([]interface{})
	var devices []client.Device
	for _, d := range configDevices {
		device := d.(map[string]interface{})
		devices = append(
			devices,
			client.Device{
				Name:        device["name"].(string),
				Description: device["description"].(string),
				IPv4Address: device["ipv4_address"].(string),
				IPv6Address: device["ipv6_address"].(string),
			},
		)

	}
	u := client.User{
		Username:  username,
		Role:      role,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		GroupId:   groupId,
		Devices:   devices,
	}
	user, err := c.CreateUser(u)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(user.Id)
	return append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "The user's role cannot be changed using the code.",
		Detail:   "There is a bug in OpenVPN Cloud API that prevents setting the user's role during the creation. All users are created as Members by default. Once it's fixed, the provider will be updated accordingly.",
	})
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	userId := d.Id()
	u, err := c.GetUserById(userId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if u == nil {
		d.SetId("")
	} else {
		d.Set("username", u.Username)
		d.Set("email", u.Email)
		d.Set("first_name", u.FirstName)
		d.Set("last_name", u.LastName)
		d.Set("group_id", u.GroupId)
		d.Set("devices", u.Devices)
	}
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	userId := d.Id()
	err := c.DeleteUser(userId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
