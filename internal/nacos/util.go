package nacos

import (
	"fmt"
	"strings"

	nacos "gitlab.zalopay.vn/top/cicd/terraform-provider-nacos/pkg/client"
)

const (
	ConfigurationIdSeparator = "/"
)

func convToConfigurationId(resourceId string) (*nacos.ConfigurationId, error) {
	parts := strings.Split(resourceId, ConfigurationIdSeparator)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid resourceId: %s", resourceId)
	}
	return &nacos.ConfigurationId{
		Namespace: parts[0],
		Group:     parts[1],
		Key:       parts[2],
	}, nil
}

func convToResourceId(namespace, group, key string) string {
	return strings.Join([]string{namespace, group, key}, ConfigurationIdSeparator)
}
