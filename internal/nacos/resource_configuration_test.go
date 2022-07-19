package nacos

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	nacos "github.com/zalopay-oss/terraform-provider-nacos/pkg/client"
)

const (
	_namespace1 = "sandbox_1"
	_namespace2 = "sandbox_2"
	_group1     = "SECRET_1"
	_group2     = "SECRET_2"

	_value1 = "test value"
	_value2 = "test value changed"
)

func TestAccNacosConfiguration_basic(t *testing.T) {
	var configuration nacos.Configuration
	rKey := fmt.Sprintf("config-key-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccNacosConfigurationPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckNacosConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNacosConfigurationConfig(rKey, nacos.Configuration{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNacosConfigurationExists("sample", &configuration),
					testAccCheckNacosConfigurationAttributes(&configuration, &nacos.Configuration{
						Namespace: _namespace1,
						Group:     _group1,
						Key:       rKey,
						Value:     _value1,
					}),
				),
			},
			// update group, re-create
			{
				Config: testAccNacosConfigurationConfig(rKey, nacos.Configuration{
					Group: _group2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNacosConfigurationExists("sample", &configuration),
					testAccCheckNacosConfigurationAttributes(&configuration, &nacos.Configuration{
						Group: _group2,
					}),
				),
			},
			// update namespace, re-create
			{
				Config: testAccNacosConfigurationConfig(rKey, nacos.Configuration{
					Namespace: _namespace2,
					Group:     _group2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNacosConfigurationExists("sample", &configuration),
					testAccCheckNacosConfigurationAttributes(&configuration, &nacos.Configuration{
						Namespace: _namespace2,
						Group:     _group2,
					}),
				),
			},
			// update value, description
			{
				Config: testAccNacosConfigurationConfig(rKey, nacos.Configuration{
					Namespace:   _namespace2,
					Group:       _group2,
					Value:       _value2,
					Description: "some description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNacosConfigurationExists("sample", &configuration),
					testAccCheckNacosConfigurationAttributes(&configuration, &nacos.Configuration{
						Value:       _value2,
						Description: "some description",
					}),
				),
			},
		},
	})
}

// testAccNacosConfigurationConfig: generate a terraform config for nacos_configuration
func testAccNacosConfigurationConfig(rName string, opts nacos.Configuration) string {
	if opts.Namespace == "" {
		opts.Namespace = _namespace1
	}
	if opts.Group == "" {
		opts.Group = _group1
	}
	if opts.Value == "" {
		opts.Value = _value1
	}

	return fmt.Sprintf(`
	resource "nacos_configuration" "sample" {
		namespace = "%s"
		group = "%s"
		key = "%s"
		value = "%s"
		description = "%s"
	}
	`, opts.Namespace,
		opts.Group,
		rName,
		opts.Value,
		opts.Description)
}

// test hooks
func testAccNacosConfigurationPreCheck(t *testing.T) {
	var missingEnvs []string
	for _, env := range []string{
		"NACOS_USERNAME",
		"NACOS_PASSWORD",
		"NACOS_ADDRESS",
	} {
		if os.Getenv(env) == "" {
			missingEnvs = append(missingEnvs, env)
		}
	}

	if len(missingEnvs) > 0 {
		t.Fatalf("%v missing for acceptance tests", missingEnvs)
	}
}

func testAccCheckNacosConfigurationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nacos_configuration" {
			continue
		}

		namespace := rs.Primary.Attributes["namespace"]
		group := rs.Primary.Attributes["group"]
		key := rs.Primary.Attributes["key"]

		c, err := testNacosClient.GetConfiguration(
			context.Background(),
			&nacos.ConfigurationId{Namespace: namespace, Group: group, Key: key})
		if err == nil {
			if c != nil {
				return fmt.Errorf("configuration (%s, %s, %s) still exists", namespace, group, key)
			}
			return nil
		}

		if !strings.Contains(err.Error(), "not found configuration") {
			return err
		}
	}

	return nil
}

// test check funcs
func testAccCheckNacosConfigurationExists(resourceName string, c *nacos.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[fmt.Sprintf("nacos_configuration.%s", resourceName)]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		configurationId, err := convToConfigurationId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error in splitting configuration ID: %v", err)
		}

		configuration, err := testNacosClient.GetConfiguration(context.Background(), configurationId)
		if err != nil {
			return err
		}

		*c = *configuration
		return nil
	}
}

func testAccCheckNacosConfigurationAttributes(configuration, want *nacos.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if want.Namespace != "" {
			if want.Namespace != configuration.Namespace {
				return fmt.Errorf("got namespace %s, want %s", configuration.Namespace, want.Namespace)
			}
		}

		if want.Group != "" {
			if want.Group != configuration.Group {
				return fmt.Errorf("got group %s, want %s", configuration.Group, want.Group)
			}
		}

		if want.Key != "" {
			if want.Key != configuration.Key {
				return fmt.Errorf("got key %s, want %s", configuration.Key, want.Key)
			}
		}

		if want.Value != "" {
			if want.Value != configuration.Value {
				return fmt.Errorf("got value %s, want %s", configuration.Value, want.Value)
			}
		}

		if want.Description != configuration.Description {
			return fmt.Errorf("got description %s, want %s", configuration.Description, want.Description)
		}

		return nil
	}
}
