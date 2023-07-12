# 微服务 用户模块

### 介绍
golang web项目
```
iris        框架
gorm        数据库
logrus      日志库
viper       配置文件读取
go-redis    redis库
```

### 项目结构
```
|-- app             //应用程序目录
|  |-- controller   //控制器文件目录
|  |-- lib          //自定义类库
|  |-- model        //数据模型定义目录
|  |-- service      //逻辑类目录
|-- config          //配置文件目录                 
|-- go.mod          //go mod依赖配置文件
|-- main.go         //主程序入口文件
|-- README.md
```


### 运行
```
    go run main.go
```

