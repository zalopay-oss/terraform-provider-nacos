package nacos

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	nacos "gitlab.zalopay.vn/top/cicd/terraform-provider-nacos/pkg/client"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"nacos": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

var testNacosConfig = nacos.Config{
	Address:     os.Getenv("NACOS_ADDRESS"),
	Username:    os.Getenv("NACOS_USERNAME"),
	Password:    os.Getenv("NACOS_PASSWORD"),
	ContextPath: os.Getenv("NACOS_CONTEXT_PATH"),
}

var testNacosClient *nacos.Client

func init() {
	if b, err := strconv.ParseBool(os.Getenv("TF_ACC")); err != nil || !b {
		return
	}

	client, err := nacos.NewClient(&testNacosConfig)
	if err != nil {
		panic("failed to create test nacos client: " + err.Error())
	}

	testNacosClient = client
}

func TestProvider(t *testing.T) {
	t.Parallel()

	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
