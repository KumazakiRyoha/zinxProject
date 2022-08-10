package main

import "github.com/KumazakiRyoha/zinxProject/znet"

/**
基于Zinx框架来开发的服务器端2应用程序
*/

func main() {
	// 创建Server实例，使用Zinx的api
	server := znet.NewServer("[zinx V0.1]")
	// 启动server
	server.Serve()
}
