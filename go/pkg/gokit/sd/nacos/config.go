package nacos

import "github.com/nacos-group/nacos-sdk-go/v2/common/constant"

type Config struct {
	//ServerConfig constant.ServerConfig `json:"serverConfig"`
	ClientConfig *constant.ClientConfig `json:"clientConfig"`
}
