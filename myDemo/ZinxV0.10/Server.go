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
	fmt.Println("recv from client msgId = ", request.GetMsgId(), ", Data = ", string(request.GetData()))
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

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====>DoConnnection is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}

	// 设置连接属性高
	fmt.Println("Set conn name,etc....")
	conn.SetProperty("Name", "雨宮ひかり")
	conn.SetProperty("Gtihub", "https://github.com/KumazakiRyoha")
	conn.SetProperty("Home", "https://github.com/KumazakiRyoha")
	conn.SetProperty("School", "明和第一高校")
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("====>DoConnection is Called...")
	fmt.Println("conn id = ", conn.GetConnID(), " is Lost...")
	// 获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if github, err := conn.GetProperty("Gtihub"); err == nil {
		fmt.Println("Gtihub = ", github)
	}
	if school, err := conn.GetProperty("School"); err == nil {
		fmt.Println("School = ", school)
	}

}

func main() {
	// 创建Server实例，使用Zinx的api
	server := znet.NewServer()
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})

	// 注册钩子函数
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	// 启动server
	server.Serve()
}
