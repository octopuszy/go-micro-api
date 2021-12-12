package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/octopuszy/go-micro-api/handler"
	"github.com/octopuszy/go-micro-api/proto/userApi"
	"github.com/octopuszy/go-micro-user/proto/user"
	util "github.com/octopuszy/micro-util"
	"github.com/opentracing/opentracing-go"
	"log"
	"net"
	"net/http"
)

func main() {
	// 注册中心
	registre := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"http://127.0.0.1:2379"}
	})

	// 链路追踪
	trancer, i, err := util.NewTrancer("test.client", "localhost:6831")
	if err != nil {
		return
	}
	defer i.Close()
	opentracing.SetGlobalTracer(trancer)

	// 熔断器面板，若需要可视化，则需要启动流服务推送上去，若仅使用降级功能不需要可视化，则不需要下面代码。
	hystrixHandler := hystrix.NewStreamHandler()
	hystrixHandler.Start()
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0","8889"), hystrixHandler)
		if err != nil {
			log.Println(err)
		}
	}()

	// 注入并开启 micro service
	ser := micro.NewService(
		micro.Name("test.api"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8086"),
		micro.Registry(registre),		// 注册中心
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),	// 链路追踪客户端
		micro.WrapClient(util.NewHystrixClientWrap()),  // 熔断

	)
	ser.Init()

	userService := user.NewUserService("test.server.user", ser.Client())

	err = userApi.RegisterUserApiHandler(ser.Server(), &handler.UserApi{UserService: userService})
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := ser.Run(); err != nil {
		log.Fatal(err);
	}
}
