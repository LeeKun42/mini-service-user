package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"iris-app/app/lib/nacos"
	"os"
	"strings"
	"time"
)

func Init() {
	loadRemoteConfig()
}

// resetLocalConfig 重置本地配置文件
func resetLocalConfig() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.Reset()
	//设置配置文件目录
	viper.AddConfigPath(path)
	//设置配置文件名称
	viper.SetConfigName("env")
	//设置配置文件类型
	viper.SetConfigType("yaml")
	//设置监听配置文件修改变化
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// loadRemoteConfig 加载nacos远程文件
func loadRemoteConfig() {
	//先加载本地配置文件
	resetLocalConfig()
	var config nacos.Config
	viper.UnmarshalKey("nacos", &config)
	var dataId string = config.DataId
	var group string = config.Group

	//获取nacos配置客户端
	client := nacos.NewConfigClient()

	//读取nacos配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		fmt.Println("nacos读取配置错误", err)
	}

	//把nacos远程配置和本地配置合并
	err = viper.MergeConfig(strings.NewReader(content))
	if err != nil {
		fmt.Println("viper解析配置失败", err)
	}

	//监听配置变化
	go func() {
		time.Sleep(time.Second * 10)
		client.ListenConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				//重置本地配置
				resetLocalConfig()
				//把nacos远程配置和本地配置合并
				err = viper.MergeConfig(strings.NewReader(data))
				if err != nil {
					fmt.Println("viper解析配置失败", err)
				}
			},
		})
	}()
}
