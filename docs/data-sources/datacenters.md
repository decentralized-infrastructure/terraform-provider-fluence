---
page_title: "fluence_datacenters Data Source - terraform-provider-fluence"
subcategory: ""
description: |-
  Fetch list of registered datacenters
---

# fluence_datacenters (Data Source)

Fetch list of registered datacenters



## Schema

### Read-Only

- `datacenters` (Attributes List) List of available datacenters (see [below for nested schema](#nestedatt--datacenters))

<a id="nestedatt--datacenters"></a>
### Nested Schema for `datacenters`

Read-Only:

- `certifications` (List of String) List of datacenter certifications
- `city_code` (String) City code
- `country_code` (String) Country code
- `id` (String) Datacenter ID
- `index` (Number) Datacenter index
- `slug` (String) Datacenter slug identifier
- `tier` (Number) Datacenter tier
