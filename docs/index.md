---
page_title: "nacos Provider"
subcategory: ""
description: |-
  Terraform provider for interacting with Nacos API.
---

# Nacos Provider
This is a Nacos Provider which provides resources to interact with [Nacos](https://github.com/alibaba/nacos) by using  [Nacos Open API](https://nacos.io/en-us/docs/open-api.html).

## Example Usage
```terraform
provider "nacos" {
  username = "username"
  password = "password"
  address = "https://nacos.example.com"
  context_path = "nacos"
}
```

## Schema
- `address` (String,required) can be set with env `NACOS_ADDRESS`, must contain protocol scheme: `https://` or `http://`
- `username` (String,required) can be set with env `NACOS_USERNAME`
- `password` (String, required) can be set with env `NACOS_PASSWORD`

### Optional
- `context_path` (String) can be set with env `NACOS_PASSWORD`, default is `nacos`
