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

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写
	fmt.Println("recv from client msgId = ", request.GetMsgId(),
		", Data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (p *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再回写
	fmt.Printf("recv from client msgId = ", request.GetMsgId(),
		", Data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 创建Server实例，使用Zinx的api
	server := znet.NewServer()
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})
	// 启动server
	server.Serve()
}
