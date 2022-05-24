package openvpncloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/patoarvizu/terraform-provider-openvpn-cloud/client"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteCreate,
		UpdateContext: resourceRouteUpdate,
		ReadContext:   resourceRouteRead,
		DeleteContext: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{client.RouteTypeIPV4, client.RouteTypeIPV6, client.RouteTypeDomain}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	networkItemId := d.Get("network_item_id").(string)
	routeType := d.Get("type").(string)
	routeValue := d.Get("value").(string)
	routeDescription := d.Get("description").(string)
	r := client.Route{
		Type:        routeType,
		Value:       routeValue,
		Description: routeDescription,
	}
	tflog.Info(ctx, "Creating OpenVPN route")
	tflog.Debug(ctx, fmt.Sprintf("Creating OpenVPN route %s, type %s, description %s", r.Value, r.Type, r.Description))
	route, err := c.CreateRoute(networkItemId, r)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Created OpenVPN route %s, type %s, id %s", route.Value, route.Type, route.Id))
	d.SetId(route.Id)
	if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
		d.Set("value", route.Subnet)
	} else if routeType == client.RouteTypeDomain {
		d.Set("value", route.Domain)
	}

	return resourceRouteUpdate(ctx, d, m)
}

func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	networkItemId := d.Get("network_item_id").(string)
	routeType := d.Get("type").(string)
	routeValue := d.Get("value").(string)
	routeDescription := d.Get("description").(string)
	routeId := d.Id()
	r := client.Route{
		Type:        routeType,
		Value:       routeValue,
		Description: routeDescription,
		Id:          routeId,
	}
	tflog.Debug(ctx, fmt.Sprintf("Updating OpenVPN route %s, type %s", r.Value, r.Type))
	err := c.UpdateRoute(networkItemId, r)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return resourceRouteRead(ctx, d, m)
}

func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Reading OpenVPN route id %s", routeId))
	r, err := c.GetRouteById(routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("Read OpenVPN route id %s, type %s, value %s", r.Id, r.Type, r.Value))
	if r == nil {
		d.SetId("")
	} else {
		d.Set("type", r.Type)
		d.Set("description", r.Description)
		if r.Type == client.RouteTypeIPV4 || r.Type == client.RouteTypeIPV6 {
			d.Set("value", r.Subnet)
		} else if r.Type == client.RouteTypeDomain {
			d.Set("resourceRouteRead", r.Domain)
		}
		d.Set("network_item_id", r.NetworkItemId)
	}

	return diags
}

func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	networkItemId := d.Get("network_item_id").(string)
	tflog.Info(ctx, "Deleting OpenVPN route")
	tflog.Debug(ctx, fmt.Sprintf("Deleting OpenVPN route id %s", routeId))
	err := c.DeleteRoute(networkItemId, routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	tflog.Debug(ctx, fmt.Sprintf("Deleted OpenVPN route id %s", routeId))

	return diag.FromErr(err)
}
