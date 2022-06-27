---
page_title: "nacos_configuration Resource - terraform-provider-nacos"
subcategory: ""
description: |-
  The configuration resource allows you to CRUD a nacos configuration.
---

# Resource `nacos_configuration`
The configuration resource allows you to CRUD a nacos configuration.

## Example Usage

```terraform
resource "nacos_configuration" "sample" {
  namespace = "sandbox"
  group = "SECRET"
  key = "test_key"
  value = "test_value"
  description = "this is the description"
}

output "sample_configuration" {
  value = "${nacos_configuration.sample.key}:${nacos_configuration.sample.value}"
}
```

## Argument Reference
Nacos configurations have 2 levels of isolation: namespace and group.

Configurations in one group of a namespace must have unique key.

- `namespace` (String, ForceNew)
- `group` (String, ForceNew)
- `key` (String, ForceNew)
- `value` (String)

### Optional
- `description` (String)