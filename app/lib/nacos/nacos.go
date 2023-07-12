package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"iris-app/app/lib"
)

type Config struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataId    string `mapstructure:"data_id"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

func NewConfigClient() config_client.IConfigClient {
	var config Config
	viper.UnmarshalKey("nacos", &config)
	serverConfigs := []constant.ServerConfig{
		{IpAddr: config.Host, Port: config.Port},
	}
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.Namespace,
			LogDir:              viper.GetString("runtime.dir") + "/logs/",
			CacheDir:            viper.GetString("runtime.dir") + "/cache",
			NotLoadCacheAtStart: true,
			Username:            config.Username,
			Password:            config.Password,
		},
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		fmt.Println("nacos初始化错误", err)
	}
	return client
}

func NewNamingClient() naming_client.INamingClient {
	var config Config
	viper.UnmarshalKey("nacos", &config)
	viper.SetConfigType("yaml")
	serverConfigs := []constant.ServerConfig{
		{IpAddr: config.Host, Port: config.Port},
	}
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         config.Namespace,
			LogDir:              viper.GetString("runtime.dir") + "/logs/",
			CacheDir:            viper.GetString("runtime.dir") + "/cache",
			NotLoadCacheAtStart: true,
			Username:            config.Username,
			Password:            config.Password,
		},
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		fmt.Println("nacos初始化错误", err)
	}
	return client
}

func RegisterService(serviceName string, ip string, port int, groupName string, metaData map[string]string) {
	ncClient := NewNamingClient()
	var err error
	if ip == "" {
		ip, err = lib.GetOutBoundIP()
		if err != nil {
			fmt.Println("没有获取到本机对外ip", err)
			return
		}
	}
	ok, err := ncClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip, //服务注册本机实例的ip
		Port:        uint64(port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		ServiceName: serviceName,
		Ephemeral:   true,
		ClusterName: "",
		GroupName:   groupName,
		Metadata:    metaData,
	})
	if err != nil {
		fmt.Println("RegisterService err", err)
	}
	if !ok {
		fmt.Println("RegisterService err 2")
	} else {
		fmt.Printf("RegisterService %s in %s success\n", serviceName, groupName)
	}
}

func GetOneHealthyInstance(serviceName string, groupName string) *model.Instance {
	ins, err := NewNamingClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	})
	if err != nil {
		fmt.Println("GetOneHealthyInstance err", err)
	}
	return ins
}
