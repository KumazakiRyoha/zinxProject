package utils

import (
	"encoding/json"
	"github.com/KumazakiRyoha/zinxProject/ziface"
	"io/ioutil"
)

type GlobeObj struct {
	/**
	Server
	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	Version          string // 当前Zinx的版本号
	MaxConn          int    //当前服务器允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作Worker池的Goroutine
	MaxWorkerTaskLen uint32 // Zinx框架允许用户最多开辟多少个Worker
}

/**
定义一个全局的Globeobj
*/
var GlobleObj *GlobeObj

/**
从配置文件加载自定义参数
*/
func (g *GlobeObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将元数据解析到struct
	err = json.Unmarshal(data, &GlobleObj)
	if err != nil {
		panic(err)
	}
}

/**
提供一个init方法，初始化GlobleObj
*/
func init() {
	GlobleObj = &GlobeObj{
		Name:             "ZinxServerApp",
		Version:          "V0.7",
		TcpPort:          8999,
		Host:             "0.0.0.0", // 代表全部IP
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的任务的数量最大值
	}

	// 应该尝试从conf/zinx.json 去加载一些用户自定义的参数
	GlobleObj.Reload()
}
