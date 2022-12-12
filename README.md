# Terraform Provider for Nacos
The Terraform Nacos Provider is a plugin for Terraform that allows for the full lifecycle management of Nacos configuration.

Currently, we only support configuration management on Nacos.

However, our goal is to support other Nacos resources as well such as Namepsace, Service Discovery...  

Public provider link https://registry.terraform.io/providers/zalopay-oss/nacos/0.1.1

# How to use
Please check  Document in `/docs`

# How to develop
### Build
```bash
OS_ARCH=<your os>_<your architecture> make install
# Example
OS_ARCH=darwin_arm64 make install
```

### Run
You can run an example in the `/examples` folder at your machine

Requisites:
- terraform installed
- accessibility to a nacos server

```bash
cd examples/

# Please use your actual credentials here
export NACOS_USERNAME=<nacos username> NACOS_PASSWORD=<nacos password> NACOS_ADDRESS=<nacos address> 
# or declare them in main.tf

terraform init
terraform apply
```

### Test
Terraform provider using acceptance tests.

Please refer to https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests when you add new test cases.

To run acceptance tests at your machine, make sure
- the accessibility to the nacos server
- create two namespaces `sandbox_1` and `sandbox_2` on nacos for testing (check in `internal/nacos/resource_configuration_test.go`)
- an installed terraform. Otherwise the test will install another terraform version when running

```bash
make testacc
```
