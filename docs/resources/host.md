---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openvpncloud_host Resource - terraform-provider-openvpncloud"
subcategory: ""
description: |-
  
---

# openvpncloud_host (Resource)

Use `openvpncloud_host` to create an OpenVPN Cloud host.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **connector** (Block Set, Min: 1) The set of connectors to be associated with this host. Can be defined more than once. (see [below for nested schema](#nestedblock--connector))
- **name** (String) The display name of the host.

### Optional

- **description** (String) The description for the UI. Defaults to `Managed by Terraform`.
- **internet_access** (String) The type of internet access provided. Valid values are `BLOCKED`, `GLOBAL_INTERNET`, or `LOCAL`. Defaults to `LOCAL`.

### Read-Only

- **id** (String) The ID of this resource.
- **system_subnets** (Set of String) The IPV4 and IPV6 subnets automatically assigned to this host.

<a id="nestedblock--connector"></a>
### Nested Schema for `connector`

Required:

Required:

- **name** (String) Name of the connector associated with this host.
- **vpn_region_id** (String) The id of the region where the connector will be deployed.

Read-Only:

- **id** (String) The ID of the connector.
- **ip_v4_address** (String) The IPV4 address of the connector.
- **ip_v6_address** (String) The IPV6 address of the connector.
- **network_item_id** (String) The host id.
- **network_item_type** (String) The network object type. This typically will be set to `HOST`.

