---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openvpncloud Provider"
subcategory: ""
description: |-
---

# openvpncloud Provider

!> **Deprecation notice:** This repository has been moved to the official [OpenVPN organization](https://github.com/OpenVPN/terraform-provider-openvpn-cloud).
Start using a new provider, **this repository will no longer be developed.**

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- **base_url** (String) The base url of your OpenVPN Cloud accout.

### Optional

- **client_id** (String, Sensitive) If not provided, it will default to the value of the `OPENVPN_CLOUD_CLIENT_ID` environment variable.
- **client_secret** (String, Sensitive) If not provided, it will default to the value of the `OPENVPN_CLOUD_CLIENT_SECRET` environment variable.
