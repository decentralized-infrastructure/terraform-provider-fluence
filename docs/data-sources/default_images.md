---
page_title: "fluence_default_images Data Source - terraform-provider-fluence"
subcategory: ""
description: |-
  Fetch list of default OS images
---

# fluence_default_images (Data Source)

Fetch list of default OS images



## Schema

### Read-Only

- `images` (Attributes List) List of available default OS images (see [below for nested schema](#nestedatt--images))

<a id="nestedatt--images"></a>
### Nested Schema for `images`

Read-Only:

- `created_at` (String) Image creation timestamp
- `distribution` (String) OS distribution
- `download_url` (String) Image download URL
- `id` (String) Image ID
- `name` (String) Image name
- `slug` (String) Image slug identifier
- `updated_at` (String) Image last update timestamp
- `username` (String) Default username for the image
