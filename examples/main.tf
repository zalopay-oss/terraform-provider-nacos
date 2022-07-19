terraform {
  required_providers {
    nacos = {
      version = "0.1.0"
      source = "zalopay-oss/nacos"
    }
  }
}

provider "nacos" {
  # address = <address or set env NACOS_ADDRESS>
  # username = <username or set env NACOS_USERNAME>
  # password = <password or set env NACOS_PASSWORD>
}

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