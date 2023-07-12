package http

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	recover2 "github.com/kataras/iris/v12/middleware/recover"
	"github.com/spf13/viper"
	"iris-app/app/http/controller"
	"iris-app/app/http/middleware"
	"iris-app/app/lib/log"
	"iris-app/app/lib/nacos"
	"iris-app/app/rpc/service"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

const ServiceName = "user"

// StartWebServer 开启web服务
func StartWebServer() {
	app := iris.New()
	app.Use(recover2.New())
	//跨域中间件
	app.UseRouter(cors.New().Handler())
	app.Use(log.HttpLogHandler)

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		fmt.Println("Shutdown server")
		// 关闭所有主机
		app.Shutdown(ctx)
	})

	/** restful风格路由 */
	app.Get("/", controller.UserController.Hello)
	app.PartyFunc("/api", func(api iris.Party) {
		api.PartyFunc("/users", func(users iris.Party) {
			users.Put("", controller.UserController.Register)
			users.Get("/token", controller.UserController.Login)
			users.Get("/new-token", controller.UserController.RefreshToken).Use(middleware.JwtAuthCheck)
			users.Delete("/token", controller.UserController.Logout).Use(middleware.JwtAuthCheck)
			users.Get("/info", controller.UserController.Info).Use(middleware.JwtAuthCheck)
		})
	})

	port := viper.GetInt("server.http")

	//注册rpc服务实例到nacos
	nacos.RegisterService(ServiceName, "", port, "http-service", map[string]string{})

	/**
	开启web服务
	参数1：监听地址和端口
	参数2：允许body多次消费
	*/
	app.Run(iris.Addr(fmt.Sprintf(":%d", port)), iris.WithoutBodyConsumptionOnUnmarshal)
}

// StartRpcServer 开启rpc服务
func StartRpcServer() {
	rpcServerPort := viper.GetInt("server.rpc")
	//注册rpc服务
	rpc.RegisterName(ServiceName, new(service.UserRpcService))

	//注册rpc服务实例到nacos
	nacos.RegisterService(ServiceName, "", rpcServerPort, "rpc-service", map[string]string{})

	//开启rpc服务
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", rpcServerPort))
	if err != nil {
		fmt.Println("ListenTCP error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
