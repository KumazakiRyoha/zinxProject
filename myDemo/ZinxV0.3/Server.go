package main

import (
	"fmt"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"github.com/KumazakiRyoha/zinxProject/znet"
)

/**
基于Zinx框架来开发的服务器端应用程序
*/

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping...")
	}
}

func main() {
	// 创建Server实例，使用Zinx的api
	server := znet.NewServer("[zinx V0.3]")
	server.AddRouter(&PingRouter{})
	// 启动server
	server.Serve()
}
